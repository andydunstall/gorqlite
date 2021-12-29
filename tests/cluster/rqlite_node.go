package cluster

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Shopify/toxiproxy/v2"
	"github.com/dunstall/gorqlite"
	log "github.com/sirupsen/logrus"
)

type Proxy interface {
	Stop()
}

type RqliteNode struct {
	ID            uint32
	HTTPPort      uint16
	RaftPort      uint16
	ProxyHTTPPort uint16
	ProxyRaftPort uint16
	Dir           string

	proxies []Proxy
	proc    io.Closer
}

func NewRqliteNode(id uint32) (*RqliteNode, error) {
	wrapper := gorqlite.NewErrorWrapper("failed to create rqlite node")

	dir, err := ioutil.TempDir("", fmt.Sprintf("node-%d", id))
	if err != nil {
		return nil, wrapper.Error(err)
	}

	httpPort := uint16(4000 + id)
	raftPort := uint16(5000 + id)

	proxyHTTPPort := uint16(6000 + id)
	proxyRaftPort := uint16(7000 + id)

	httpProxy := toxiproxy.NewProxy()
	httpProxy.Listen = fmt.Sprintf("0.0.0.0:%d", proxyHTTPPort)
	httpProxy.Upstream = fmt.Sprintf("localhost:%d", httpPort)
	if err = httpProxy.Start(); err != nil {
		return &RqliteNode{}, wrapper.Error(err)
	}

	raftProxy := toxiproxy.NewProxy()
	raftProxy.Listen = fmt.Sprintf("0.0.0.0:%d", proxyRaftPort)
	raftProxy.Upstream = fmt.Sprintf("localhost:%d", raftPort)
	if err = raftProxy.Start(); err != nil {
		return &RqliteNode{}, wrapper.Error(err)
	}

	rqliteProc, err := NewRqliteProc(id, httpPort, proxyHTTPPort, raftPort, proxyRaftPort, dir)
	if err != nil {
		return &RqliteNode{}, wrapper.Error(err)
	}

	return &RqliteNode{
		ID:            id,
		HTTPPort:      httpPort,
		RaftPort:      raftPort,
		ProxyHTTPPort: proxyHTTPPort,
		ProxyRaftPort: proxyRaftPort,
		Dir:           dir,
		proxies:       []Proxy{httpProxy, raftProxy},
		proc:          &rqliteProc,
	}, nil
}

func NewRqliteNodeWithJoin(id uint32, joinPort uint16) (*RqliteNode, error) {
	wrapper := gorqlite.NewErrorWrapper("failed to create rqlite node")

	dir, err := ioutil.TempDir("", fmt.Sprintf("node-%d", id))
	if err != nil {
		return nil, wrapper.Error(err)
	}

	httpPort := uint16(4000 + id)
	raftPort := uint16(5000 + id)

	proxyHTTPPort := uint16(6000 + id)
	proxyRaftPort := uint16(7000 + id)

	httpProxy := toxiproxy.NewProxy()
	httpProxy.Listen = fmt.Sprintf("0.0.0.0:%d", proxyHTTPPort)
	httpProxy.Upstream = fmt.Sprintf("localhost:%d", httpPort)
	if err = httpProxy.Start(); err != nil {
		return &RqliteNode{}, wrapper.Error(err)
	}

	raftProxy := toxiproxy.NewProxy()
	raftProxy.Listen = fmt.Sprintf("0.0.0.0:%d", proxyRaftPort)
	raftProxy.Upstream = fmt.Sprintf("localhost:%d", raftPort)
	if err = raftProxy.Start(); err != nil {
		return &RqliteNode{}, wrapper.Error(err)
	}

	rqliteProc, err := NewRqliteProcWithJoin(id, httpPort, proxyHTTPPort, raftPort, proxyRaftPort, dir, joinPort)
	if err != nil {
		return nil, wrapper.Error(err)
	}

	return &RqliteNode{
		ID:            id,
		HTTPPort:      httpPort,
		RaftPort:      raftPort,
		ProxyHTTPPort: proxyHTTPPort,
		ProxyRaftPort: proxyRaftPort,
		Dir:           dir,
		proxies:       []Proxy{httpProxy, raftProxy},
		proc:          &rqliteProc,
	}, nil
}

func (n *RqliteNode) Reboot(duration int64, timeout bool) error {
	wrapper := gorqlite.NewErrorWrapper("failed to reboot rqlite node")

	if err := n.Close(); err != nil {
		return wrapper.Error(err)
	}

	if timeout {
		<-time.After(time.Duration(duration) * time.Millisecond)
	}

	httpProxy := toxiproxy.NewProxy()
	httpProxy.Listen = fmt.Sprintf("0.0.0.0:%d", n.ProxyHTTPPort)
	httpProxy.Upstream = fmt.Sprintf("localhost:%d", n.HTTPPort)
	if err := httpProxy.Start(); err != nil {
		return wrapper.Error(err)
	}

	raftProxy := toxiproxy.NewProxy()
	raftProxy.Listen = fmt.Sprintf("0.0.0.0:%d", n.ProxyRaftPort)
	raftProxy.Upstream = fmt.Sprintf("localhost:%d", n.RaftPort)
	if err := raftProxy.Start(); err != nil {
		return wrapper.Error(err)
	}

	rqliteProc, err := NewRqliteProc(n.ID, n.HTTPPort, n.ProxyHTTPPort, n.RaftPort, n.ProxyRaftPort, n.Dir)
	if err != nil {
		return wrapper.Error(err)
	}

	n.proxies = []Proxy{httpProxy, raftProxy}
	n.proc = &rqliteProc

	return nil
}

func (n *RqliteNode) Close() error {
	for _, p := range n.proxies {
		p.Stop()
	}
	if err := n.proc.Close(); err != nil {
		return gorqlite.WrapError(err, "failed to close node")
	}
	return nil
}

type RqliteProc struct {
	id   uint32
	proc *os.Process
}

func NewRqliteProc(id uint32, httpListenPort uint16, httpProxyPort uint16, raftListenPort uint16, raftProxyPort uint16, dir string) (RqliteProc, error) {
	wrapper := gorqlite.NewErrorWrapper("failed to create rqlite proc")

	lg, err := newLogger(fmt.Sprintf("rqlited-%d", id))
	if err != nil {
		return RqliteProc{}, wrapper.Error(err)
	}

	cmd := exec.Command(
		"rqlited",
		"-node-id",
		fmt.Sprintf("%d", id),
		"-http-addr",
		fmt.Sprintf("0.0.0.0:%d", httpListenPort),
		"-http-adv-addr",
		fmt.Sprintf("0.0.0.0:%d", httpProxyPort),
		"-raft-addr",
		fmt.Sprintf("localhost:%d", raftListenPort),
		"-raft-adv-addr",
		fmt.Sprintf("0.0.0.0:%d", raftProxyPort),
		dir,
	)
	log.WithFields(log.Fields{
		"cmd": strings.Join(cmd.Args, " "),
	}).Debug("running command")
	cmd.Stdout = lg
	cmd.Stderr = lg
	if err := cmd.Start(); err != nil {
		return RqliteProc{}, wrapper.Error(err)
	}

	log.WithFields(log.Fields{
		"id":               id,
		"http_listen_port": httpListenPort,
		"http_proxy_port":  httpProxyPort,
		"raft_listen_port": raftListenPort,
		"raft_proxy_port":  raftProxyPort,
		"log":              lg.Path,
	}).Info("started rqlited")

	return RqliteProc{
		id:   id,
		proc: cmd.Process,
	}, nil
}

func NewRqliteProcWithJoin(id uint32, httpListenPort uint16, httpProxyPort uint16, raftListenPort uint16, raftProxyPort uint16, dir string, joinPort uint16) (RqliteProc, error) {
	wrapper := gorqlite.NewErrorWrapper("failed to create rqlite proc")

	lg, err := newLogger(fmt.Sprintf("rqlited-%d", id))
	if err != nil {
		return RqliteProc{}, wrapper.Error(err)
	}

	cmd := exec.Command(
		"rqlited",
		"-node-id",
		fmt.Sprintf("%d", id),
		"-http-addr",
		fmt.Sprintf("0.0.0.0:%d", httpListenPort),
		"-http-adv-addr",
		fmt.Sprintf("0.0.0.0:%d", httpProxyPort),
		"-raft-addr",
		fmt.Sprintf("localhost:%d", raftListenPort),
		"-raft-adv-addr",
		fmt.Sprintf("0.0.0.0:%d", raftProxyPort),
		"-join",
		fmt.Sprintf("0.0.0.0:%d", joinPort),
		dir,
	)
	log.WithFields(log.Fields{
		"cmd": strings.Join(cmd.Args, " "),
	}).Debug("running command")
	cmd.Stdout = lg
	cmd.Stderr = lg
	if err := cmd.Start(); err != nil {
		return RqliteProc{}, wrapper.Error(err)
	}

	log.WithFields(log.Fields{
		"id":               id,
		"http_listen_port": httpListenPort,
		"http_proxy_port":  httpProxyPort,
		"raft_listen_port": raftListenPort,
		"raft_proxy_port":  raftProxyPort,
		"log":              lg.Path,
	}).Info("started rqlited")

	return RqliteProc{
		id:   id,
		proc: cmd.Process,
	}, nil
}

func (t *RqliteProc) Close() error {
	if err := t.proc.Kill(); err != nil {
		return gorqlite.WrapError(err, "failed to close rqlite proc")
	}
	_, err := t.proc.Wait()
	if err != nil {
		return gorqlite.WrapError(err, "failed to close rqlite proc")
	}

	log.WithFields(log.Fields{
		"id": t.id,
	}).Info("stopped rqlited")

	return nil
}

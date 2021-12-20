package cluster

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type RqliteNode struct {
	ID            uint32
	HTTPPort      uint16
	RaftPort      uint16
	ProxyHTTPPort uint16
	ProxyRaftPort uint16
	Dir           string

	nodes []io.Closer
}

func NewRqliteNode(id uint32) (*RqliteNode, error) {
	dir, err := ioutil.TempDir("", fmt.Sprintf("node-%d", id))
	if err != nil {
		return nil, fmt.Errorf("failed to create node dir: %s", err)
	}

	httpPort := uint16(4000 + id)
	raftPort := uint16(5000 + id)

	proxyHTTPPort := uint16(6000 + id)
	proxyRaftPort := uint16(7000 + id)

	nodes := []io.Closer{}

	httpProxy, err := NewToxiproxyNode(uuid.New().String(), httpPort, proxyHTTPPort)
	if err != nil {
		return &RqliteNode{}, err
	}
	nodes = append(nodes, &httpProxy)

	raftProxy, err := NewToxiproxyNode(uuid.New().String(), raftPort, proxyRaftPort)
	if err != nil {
		return &RqliteNode{}, err
	}
	nodes = append(nodes, &raftProxy)

	rqliteProc, err := NewRqliteProc(id, httpPort, proxyHTTPPort, raftPort, proxyRaftPort, dir)
	if err != nil {
		return &RqliteNode{}, err
	}
	nodes = append(nodes, &rqliteProc)

	return &RqliteNode{
		ID:            id,
		HTTPPort:      httpPort,
		RaftPort:      raftPort,
		ProxyHTTPPort: proxyHTTPPort,
		ProxyRaftPort: proxyRaftPort,
		Dir:           dir,
		nodes:         nodes,
	}, nil
}

func NewRqliteNodeWithJoin(id uint32, joinPort uint16) (*RqliteNode, error) {
	dir, err := ioutil.TempDir("", fmt.Sprintf("node-%d", id))
	if err != nil {
		return nil, fmt.Errorf("failed to create node dir: %s", err)
	}

	httpPort := uint16(4000 + id)
	raftPort := uint16(5000 + id)

	proxyHTTPPort := uint16(6000 + id)
	proxyRaftPort := uint16(7000 + id)

	nodes := []io.Closer{}

	httpProxy, err := NewToxiproxyNode(uuid.New().String(), httpPort, proxyHTTPPort)
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, &httpProxy)

	raftProxy, err := NewToxiproxyNode(uuid.New().String(), raftPort, proxyRaftPort)
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, &raftProxy)

	rqliteProc, err := NewRqliteProcWithJoin(id, httpPort, proxyHTTPPort, raftPort, proxyRaftPort, dir, joinPort)
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, &rqliteProc)

	return &RqliteNode{
		ID:            id,
		HTTPPort:      httpPort,
		RaftPort:      raftPort,
		ProxyHTTPPort: proxyHTTPPort,
		ProxyRaftPort: proxyRaftPort,
		Dir:           dir,
		nodes:         nodes,
	}, nil
}

func (n *RqliteNode) Reboot(duration int64, timeout bool) error {
	if err := n.Close(); err != nil {
		return err
	}

	if timeout {
		<-time.After(time.Duration(duration) * time.Millisecond)
	}

	nodes := []io.Closer{}
	httpProxy, err := NewToxiproxyNode(uuid.New().String(), n.HTTPPort, n.ProxyHTTPPort)
	if err != nil {
		return err
	}
	nodes = append(nodes, &httpProxy)

	raftProxy, err := NewToxiproxyNode(uuid.New().String(), n.RaftPort, n.ProxyRaftPort)
	if err != nil {
		return err
	}
	nodes = append(nodes, &raftProxy)

	rqliteProc, err := NewRqliteProc(n.ID, n.HTTPPort, n.ProxyHTTPPort, n.RaftPort, n.ProxyRaftPort, n.Dir)
	if err != nil {
		return err
	}
	nodes = append(nodes, &rqliteProc)

	n.nodes = nodes

	return nil
}

func (n *RqliteNode) Close() error {
	var err error
	for _, node := range n.nodes {
		if closeErr := node.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

type RqliteProc struct {
	id   uint32
	proc *os.Process
}

func NewRqliteProc(id uint32, httpListenPort uint16, httpProxyPort uint16, raftListenPort uint16, raftProxyPort uint16, dir string) (RqliteProc, error) {
	lg, err := newLogger(fmt.Sprintf("rqlited-%d", id))
	if err != nil {
		return RqliteProc{}, fmt.Errorf("failed to open log file: %s", err)
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
		return RqliteProc{}, err
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
	lg, err := newLogger(fmt.Sprintf("rqlited-%d", id))
	if err != nil {
		return RqliteProc{}, fmt.Errorf("failed to open log file: %s", err)
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
		return RqliteProc{}, err
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
		return err
	}
	_, err := t.proc.Wait()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id": t.id,
	}).Info("stopped rqlited")

	return nil
}

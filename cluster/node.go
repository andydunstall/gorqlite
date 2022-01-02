package cluster

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/toxiproxy/v2"
	"github.com/dunstall/gorqlite"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	id      uint32
	path    string
	conf    *nodeConfig
	proxies []*toxiproxy.Proxy
	proc    *os.Process
}

func NewNode(path string, id uint32, opts ...NodeOption) (*Node, error) {
	conf := defaultNodeConfig()
	for _, opt := range opts {
		opt(conf)
	}

	wrapper := newErrorWrapper("failed to create rqlite node")

	// Add proxies between the advertised API and actual API.
	apiProxy := toxiproxy.NewProxy()
	apiProxy.Listen = conf.APIAdvAddr
	apiProxy.Upstream = conf.APIAddr
	if err := apiProxy.Start(); err != nil {
		return &Node{}, wrapper.Error(err)
	}
	raftProxy := toxiproxy.NewProxy()
	raftProxy.Listen = conf.RaftAdvAddr
	raftProxy.Upstream = conf.RaftAddr
	if err := raftProxy.Start(); err != nil {
		return &Node{}, wrapper.Error(err)
	}

	proc, err := runNode(path, id, conf)
	if err != nil {
		return &Node{}, wrapper.Error(err)
	}

	return &Node{
		path:    path,
		id:      id,
		conf:    conf,
		proxies: []*toxiproxy.Proxy{apiProxy, raftProxy},
		proc:    proc,
	}, nil
}

func (n *Node) ID() uint32 {
	return n.id
}

func (n *Node) APIAdvAddr() string {
	return n.conf.APIAdvAddr
}

func (n *Node) RaftAdvAddr() string {
	return n.conf.RaftAdvAddr
}

func (n *Node) Reboot(duration int64, timeout bool) error {
	wrapper := newErrorWrapper("failed to reboot rqlite node")

	if err := n.Close(); err != nil {
		return wrapper.Error(err)
	}

	if timeout {
		<-time.After(time.Duration(duration) * time.Millisecond)
	}

	// Add proxies between the advertised API and actual API.
	apiProxy := toxiproxy.NewProxy()
	apiProxy.Listen = n.conf.APIAdvAddr
	apiProxy.Upstream = n.conf.APIAddr
	if err := apiProxy.Start(); err != nil {
		return wrapper.Error(err)
	}
	raftProxy := toxiproxy.NewProxy()
	raftProxy.Listen = n.conf.RaftAdvAddr
	raftProxy.Upstream = n.conf.RaftAddr
	if err := raftProxy.Start(); err != nil {
		return wrapper.Error(err)
	}

	proc, err := runNode(n.path, n.id, n.conf)
	if err != nil {
		return wrapper.Error(err)
	}

	n.proxies = []*toxiproxy.Proxy{apiProxy, raftProxy}
	n.proc = proc

	return nil
}

func (n *Node) WaitForLeader(ctx context.Context) (string, error) {
	ticker := time.NewTicker(250 * time.Millisecond)
	var err error
	for {
		select {
		case <-ctx.Done():
			if err != nil {
				return "", wrapError(err, "failed to get leader")
			}
			return "", newError("failed to get leader: timed out")
		case <-ticker.C:
			var leader string
			if leader, err = n.leader(ctx); err == nil {
				return leader, nil
			}
		}
	}
}

// WaitForAllFSM waits until all outstanding database commands have actually
// been applied to the database i.e. state machine.
func (n *Node) WaitForAllFSM(ctx context.Context) (int, error) {
	ticker := time.NewTicker(250 * time.Millisecond)
	var err error
	for {
		select {
		case <-ctx.Done():
			if err != nil {
				return 0, err
			}
			return 0, newError("failed to wait for all fsm: timed out")
		case <-ticker.C:
			var status gorqlite.Status
			status, err = n.Status(ctx)
			if err != nil {
				continue
			}

			if status.Store.FSMIndex != status.Store.DBAppliedIndex {
				err = newError("fsmIndex != dbAppliedIndex (%d != %d)", status.Store.FSMIndex, status.Store.DBAppliedIndex)
				continue
			}
			return status.Store.FSMIndex, nil
		}
	}
}

func (n *Node) Status(ctx context.Context) (gorqlite.Status, error) {
	conn := gorqlite.Open([]string{n.APIAdvAddr()})
	status, err := conn.StatusWithContext(ctx)
	if err != nil {
		return status, wrapError(err, "failed to get status")
	}
	return status, nil
}

func (n *Node) Close() error {
	for _, p := range n.proxies {
		p.Stop()
	}

	if err := n.proc.Kill(); err != nil {
		return wrapError(err, "failed to close rqlite proc")
	}
	_, err := n.proc.Wait()
	if err != nil {
		return wrapError(err, "failed to close rqlite proc")
	}

	log.WithFields(log.Fields{
		"id": n.id,
	}).Info("stopped rqlited")

	return nil
}

func (n *Node) leader(ctx context.Context) (string, error) {
	conn := gorqlite.Open([]string{n.APIAdvAddr()})
	status, err := conn.StatusWithContext(ctx)
	if err != nil {
		return "", wrapError(err, "failed to get leader")
	}
	return status.Store.Leader.Addr, nil
}

func runNode(path string, id uint32, conf *nodeConfig) (*os.Process, error) {
	wrapper := newErrorWrapper("failed to create rqlite proc")

	lg, err := newLogger(fmt.Sprintf("rqlited-%d", id))
	if err != nil {
		return nil, wrapper.Error(err)
	}

	args := []string{
		"-node-id",
		fmt.Sprintf("%d", id),
		"-http-addr",
		conf.APIAddr,
		"-http-adv-addr",
		conf.APIAdvAddr,
		"-raft-addr",
		conf.RaftAddr,
		"-raft-adv-addr",
		conf.RaftAdvAddr,
		"-raft-snap",
		strconv.Itoa(conf.RaftSnapThreshold),
		"-raft-snap-int",
		conf.RaftSnapInterval,
	}
	if conf.Join != "" {
		args = append(args, "-join")
		args = append(args, conf.Join)
	}
	if !conf.RaftVoter {
		args = append(args, "-raft-non-voter")
	}

	args = append(args, conf.Dir)
	cmd := exec.Command(path, args...)
	log.WithFields(log.Fields{
		"cmd": strings.Join(cmd.Args, " "),
	}).Debug("running command")
	cmd.Stdout = lg
	cmd.Stderr = lg
	if err := cmd.Start(); err != nil {
		return nil, wrapper.Error(err)
	}

	log.WithFields(log.Fields{
		"id":            id,
		"api_addr":      conf.APIAddr,
		"api_adv_addr":  conf.APIAdvAddr,
		"raft_addr":     conf.RaftAddr,
		"raft_adv_addr": conf.RaftAdvAddr,
		"log":           lg.Path,
		"join":          conf.Join,
	}).Info("started rqlited")

	return cmd.Process, nil
}

package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/dunstall/gorqlite"
	log "github.com/sirupsen/logrus"
)

type Cluster struct {
	nodes map[uint32]*RqliteNode
}

func NewCluster() *Cluster {
	return &Cluster{
		nodes: map[uint32]*RqliteNode{},
	}
}

// OpenCluster creates a new cluster with `numNodes` nodes where the first
// node is the leader that the other nodes join.
func OpenCluster(numNodes uint32) (*Cluster, error) {
	cluster := NewCluster()

	var leaderPort uint16
	for id := uint32(1); id <= numNodes; id++ {
		if leaderPort == 0 {
			node, err := NewRqliteNode(id)
			if err != nil {
				return nil, err
			}
			cluster.AddNode(id, node)
			leaderPort = node.ProxyHTTPPort
		} else {
			node, err := NewRqliteNodeWithJoin(id, leaderPort)
			if err != nil {
				return nil, err
			}
			cluster.AddNode(id, node)
		}
	}

	return cluster, nil
}

// WaitForHealthy waits for all nodes to return a status with a consistent
// view of the clusters nodes and leader.
func (c *Cluster) WaitForHealthy(ctx context.Context) bool {
	ticker := time.NewTicker(250 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if c.isHealthy() {
				return true
			}
		}
	}
}

func (c *Cluster) NodeAddrs() map[uint32]string {
	nodeAddresses := make(map[uint32]string)
	for id, node := range c.nodes {
		addr := fmt.Sprintf("localhost:%d", node.ProxyHTTPPort)
		nodeAddresses[id] = addr
	}
	return nodeAddresses
}

func (c *Cluster) Addrs() []string {
	addrs := []string{}
	for _, addr := range c.NodeAddrs() {
		addrs = append(addrs, addr)
	}
	return addrs
}

func (c *Cluster) RandomNodeAddr() string {
	nodes := c.NodeAddrs()
	if len(nodes) == 0 {
		return ""
	}

	ids := []uint32{}
	for id := range nodes {
		ids = append(ids, id)
	}
	return nodes[ids[rand.Int()%len(ids)]]
}

func (c *Cluster) AddNode(id uint32, node *RqliteNode) {
	c.nodes[id] = node
}

func (c *Cluster) RemoveNode(id uint32) {
	if node, ok := c.nodes[id]; ok {
		node.Close()
		delete(c.nodes, id)
	}
}

func (c *Cluster) Close() error {
	if c.nodes == nil {
		return nil
	}

	var err error
	for _, proc := range c.nodes {
		if closeErr := proc.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

func (c *Cluster) isHealthy() bool {
	var knownLeader gorqlite.Leader
	for id, addr := range c.NodeAddrs() {
		lg := log.WithFields(log.Fields{
			"node_id":   id,
			"node_addr": addr,
		})

		statusClient := gorqlite.NewStatusAPIClient(gorqlite.NewHTTPAPIClient([]string{addr}))
		status, err := statusClient.Status()
		if err != nil {
			lg.Debugf("failed to get status: %s", err)
			return false
		}

		if fmt.Sprintf("%d", id) != status.Store.NodeID {
			lg.Debugf("node returned incorrect id: %s", status.Store.NodeID)
			return false
		}

		// Check all nodes agree on the leader.
		if status.Store.Leader.NodeID == "" {
			lg.Debug("node has no leader")
			return false
		}
		if knownLeader.NodeID == "" {
			knownLeader = status.Store.Leader
		}
		if status.Store.Leader != knownLeader {
			lg.Debugf("nodes have inconsistent leaders: %#v != %#v", status.Store.Leader, knownLeader)
			return false
		}

		// Check nodes have discovered each other.
		if len(status.Store.Nodes) != len(c.nodes) {
			lg.Debug("node has not discovered cluster nodes")
			return false
		}

		lg.Debug("node healthy")
	}
	return true
}

package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/dunstall/gorqlite"
	log "github.com/sirupsen/logrus"
)

type Cluster struct {
	nodes map[uint32]*Node
}

func NewCluster() *Cluster {
	return &Cluster{
		nodes: map[uint32]*Node{},
	}
}

// OpenCluster creates a new cluster with `numNodes` nodes where the first
// node is the leader that the other nodes join.
func OpenCluster(numNodes uint32) (*Cluster, error) {
	cluster := NewCluster()

	var leaderAddr string
	for id := uint32(1); id <= numNodes; id++ {
		if leaderAddr == "" {
			node, err := NewNode("rqlited", id)
			if err != nil {
				return nil, err
			}
			cluster.AddNode(node)
			leaderAddr = node.APIAdvAddr()
		} else {
			node, err := NewNode("rqlited", id, WithJoin(leaderAddr))
			if err != nil {
				return nil, err
			}
			cluster.AddNode(node)
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
		nodeAddresses[id] = node.APIAdvAddr()
	}
	return nodeAddresses
}

func (c *Cluster) NodeRaftAddrs() map[uint32]string {
	nodeAddresses := make(map[uint32]string)
	for id, node := range c.nodes {
		nodeAddresses[id] = node.RaftAdvAddr()
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

func (c *Cluster) AddNode(node *Node) {
	c.nodes[node.ID()] = node
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
	var knownLeader gorqlite.LeaderInfo
	for id, addr := range c.NodeAddrs() {
		lg := log.WithFields(log.Fields{
			"node_id":   id,
			"node_addr": addr,
		})

		conn := gorqlite.Open([]string{addr})
		status, err := conn.Status()
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

// RunDefaultCluster runs a cluster with 3 nodes and waits for it to be
// healthy.
func RunDefaultCluster() (*Cluster, error) {
	cluster, err := OpenCluster(3)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if !cluster.WaitForHealthy(ctx) {
		return nil, newError("timed out waiting for healthy")
	}

	return cluster, nil
}

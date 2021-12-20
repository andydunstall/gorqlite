package cluster

type Cluster struct {
	nodes map[uint32]Node
}

func NewCluster() Cluster {
	return Cluster{
		nodes: map[uint32]Node{},
	}
}

func (c *Cluster) AddNode(id uint32, node Node) {
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

//go:build system

package tests

import (
	"fmt"
	"testing"

	"github.com/dunstall/gorqlite/tests/cluster"
	"github.com/stretchr/testify/require"
)

func TestSimple_PutThenGet(t *testing.T) {
	require := require.New(t)

	c := cluster.NewCluster()
	defer c.Close()

	numNodes := uint32(3)
	nodeAddresses := []string{}
	for id := uint32(1); id <= numNodes; id += 1 {
		node, err := cluster.NewRqliteNode(id)
		require.Nil(err, "failed to start node")
		c.AddNode(id, node)

		addr := fmt.Sprintf("http://localhost:%d", node.ProxyHTTPPort)
		nodeAddresses = append(nodeAddresses, addr)
	}

	// TODO(AD)
	fmt.Println(nodeAddresses)
}

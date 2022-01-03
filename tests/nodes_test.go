//go:build system

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/dunstall/gorqlite"
	"github.com/dunstall/gorqlite/cluster"
	"github.com/stretchr/testify/require"
)

func TestNodesAPIClient_PeerStatus(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	for _, addr := range cluster.NodeAddrs() {
		conn := gorqlite.Open([]string{addr})
		nodes, err := conn.Nodes()
		require.Nil(err)

		require.Equal(3, len(nodes))
		for _, n := range nodes {
			require.True(n.Reachable)
		}
	}
}

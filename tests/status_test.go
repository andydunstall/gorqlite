//go:build system

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dunstall/gorqlite"
	"github.com/dunstall/gorqlite/tests/cluster"
	"github.com/stretchr/testify/require"
)

func TestStatusAPIClient_PeerStatus(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	expectedLeader := gorqlite.Leader{
		Addr:   "0.0.0.0:7001",
		NodeID: "1",
	}
	expectedNodes := []gorqlite.Node{
		gorqlite.Node{Addr: "0.0.0.0:7001", ID: "1", Suffrage: "Voter"},
		gorqlite.Node{Addr: "0.0.0.0:7002", ID: "2", Suffrage: "Voter"},
		gorqlite.Node{Addr: "0.0.0.0:7003", ID: "3", Suffrage: "Voter"},
	}

	for id, addr := range cluster.NodeAddrs() {
		statusClient := gorqlite.NewStatusAPIClient(gorqlite.NewHTTPAPIClient([]string{addr}))
		status, err := statusClient.Status()
		require.Nil(err)

		require.Equal(fmt.Sprintf("%d", id), status.Store.NodeID)
		require.Equal(expectedLeader, status.Store.Leader)
		require.Equal(expectedNodes, status.Store.Nodes)
	}
}

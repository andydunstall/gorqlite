//go:build system

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dunstall/gorqlite"
	"github.com/dunstall/gorqlite/cluster"
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

	raftAddrs := cluster.NodeRaftAddrs()

	expectedLeader := gorqlite.LeaderInfo{
		Addr:   raftAddrs[1],
		NodeID: "1",
	}
	expectedNodes := []gorqlite.NodeInfo{
		{Addr: raftAddrs[1], ID: "1", Suffrage: "Voter"},
		{Addr: raftAddrs[2], ID: "2", Suffrage: "Voter"},
		{Addr: raftAddrs[3], ID: "3", Suffrage: "Voter"},
	}

	for id, addr := range cluster.NodeAddrs() {
		conn := gorqlite.Connect([]string{addr})
		status, err := conn.Status()
		require.Nil(err)

		require.Equal(fmt.Sprintf("%d", id), status.Store.NodeID)
		require.Equal(expectedLeader, status.Store.Leader)
		require.Equal(expectedNodes, status.Store.Nodes)
	}
}

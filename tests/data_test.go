//go:build system

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/dunstall/gorqlite"
	"github.com/dunstall/gorqlite/tests/cluster"
	"github.com/stretchr/testify/require"
)

func TestDataAPIClient_QueryXXX(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	dataClient := gorqlite.NewDataAPIClient(cluster.RandomNodeAddr())
	err = dataClient.Query([]string{"CREATE TABLE foo (id integer not null primary key, bar text)"})

	require.Nil(err)
}

func TestDataAPIClient_ExecuteXXX(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	dataClient := gorqlite.NewDataAPIClient(cluster.RandomNodeAddr())
	err = dataClient.Execute([]string{"SELECT * FROM foo"})

	require.Nil(err)
}

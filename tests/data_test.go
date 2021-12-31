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

func TestDataAPIClient_ExecuteThenQueryResults(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	conn := gorqlite.Open(cluster.Addrs())

	// Create table.
	execResults, err := conn.Execute([]string{
		"CREATE TABLE foo (id integer not null primary key, name text)",
	})
	require.Nil(err)
	require.Equal("", execResults.GetFirstError())
	require.Equal(1, len(execResults.Results))

	// Insert one row.
	execResults, err = conn.Execute([]string{
		`INSERT INTO foo(name) VALUES("fiona")`,
	})
	require.Nil(err)
	require.Equal("", execResults.GetFirstError())
	require.Equal(1, len(execResults.Results))
	require.Equal(int64(1), execResults.Results[0].RowsAffected)
	require.Equal(int64(1), execResults.Results[0].LastInsertId)

	queryResults, err := conn.Query([]string{
		`SELECT name FROM foo WHERE name="fiona"`,
	})
	require.Nil(err)
	require.Equal("", queryResults.GetFirstError())
	require.Equal(gorqlite.QueryRows{
		Columns: []string{"name"},
		Types:   []string{"text"},
		Values:  [][]interface{}{{"fiona"}},
	}, queryResults.Results[0])

	// Update one row.
	execResults, err = conn.Execute([]string{
		`UPDATE foo SET name="justin" WHERE name="fiona"`,
	})
	require.Nil(err)
	require.Equal("", execResults.GetFirstError())
	require.Equal(1, len(execResults.Results))
	require.Equal(int64(1), execResults.Results[0].RowsAffected)

	queryResults, err = conn.Query([]string{
		`SELECT name FROM foo WHERE name="justin"`,
	})
	require.Nil(err)
	require.Equal("", queryResults.GetFirstError())
	require.Equal(gorqlite.QueryRows{
		Columns: []string{"name"},
		Types:   []string{"text"},
		Values:  [][]interface{}{{"justin"}},
	}, queryResults.Results[0])

	// Delete one row.
	execResults, err = conn.Execute([]string{
		`DELETE FROM foo WHERE name="justin"`,
	})
	require.Nil(err)
	require.Equal("", execResults.GetFirstError())
	require.Equal(1, len(execResults.Results))
	require.Equal(int64(1), execResults.Results[0].RowsAffected)

	queryResults, err = conn.Query([]string{
		`SELECT COUNT(id) AS idCount FROM foo`,
	})
	require.Nil(err)
	require.Equal("", queryResults.GetFirstError())
	require.Equal(gorqlite.QueryRows{
		Columns: []string{"idCount"},
		Types:   []string{""},
		Values:  [][]interface{}{{float64(0)}},
	}, queryResults.Results[0])

	// Insert multiple rows.
	sql := []string{}
	numRows := 100
	for i := 0; i < numRows; i++ {
		sql = append(sql, fmt.Sprintf(`INSERT INTO foo(name) VALUES("justin-%d")`, i))
	}
	execResults, err = conn.Execute(sql, gorqlite.WithTransaction(true))
	require.Nil(err)
	require.Equal("", execResults.GetFirstError())
	require.Equal(numRows, len(execResults.Results))
	for i := 0; i < numRows; i++ {
		result := execResults.Results[i]
		require.Equal(int64(1), result.RowsAffected)
		require.Equal(int64(i+1), result.LastInsertId)
	}

	queryResults, err = conn.Query([]string{
		`SELECT COUNT(*) AS total FROM foo WHERE name like("justin-%")`,
	}, gorqlite.WithConsistency("strong"))
	require.Nil(err)
	require.Equal("", queryResults.GetFirstError())
	require.Equal(gorqlite.QueryRows{
		Columns: []string{"total"},
		Types:   []string{""},
		Values:  [][]interface{}{{float64(numRows)}},
	}, queryResults.Results[0])

	sql = []string{}
	for i := 0; i < numRows; i++ {
		sql = append(sql, fmt.Sprintf(`SELECT name FROM foo WHERE name="justin-%d"`, i))
	}
	queryResults, err = conn.Query(sql)
	require.Nil(err)
	require.Equal("", queryResults.GetFirstError())
	require.Equal(numRows, len(queryResults.Results))
	for i, result := range queryResults.Results {
		require.Equal(gorqlite.QueryRows{
			Columns: []string{"name"},
			Types:   []string{"text"},
			Values:  [][]interface{}{{fmt.Sprintf("justin-%d", i)}},
		}, result)
	}

	// Drop the table.
	execResults, err = conn.Execute([]string{
		`DROP TABLE foo`,
	})
	require.Nil(err)
	require.Equal("", execResults.GetFirstError())
	require.Equal(1, len(execResults.Results))
	require.Equal(int64(1), execResults.Results[0].RowsAffected)
}

func TestDataAPIClient_QueryInvalidCommand(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	conn := gorqlite.Open(cluster.Addrs())
	_, err = conn.Execute([]string{
		"CREATE TABLE foo (id integer not null primary key, bar text)",
		"INSERT INTO foo (bar) values ('baz')",
	})
	require.Nil(err)
	result, err := conn.Query([]string{
		"SELECT * FROM foo",
		// Table does not exist.
		"SELECT * FROM bar",
	})

	require.Nil(err)
	require.Equal("no such table: bar", result.GetFirstError())
}

func TestDataAPIClient_ExecuteInvalidCommand(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	conn := gorqlite.Open(cluster.Addrs())
	result, err := conn.Execute([]string{
		"CREATE TABLE foo (id integer not null primary key, bar text)",
		// Table does not exist.
		"INSERT INTO baz (bar) values ('bar')",
	})

	require.Nil(err)

	expectedResult := gorqlite.ExecuteResponse{
		Results: []gorqlite.ExecuteResult{
			{
				LastInsertId: 0,
				RowsAffected: 0,
				Error:        "",
				Time:         0,
			},
			{
				LastInsertId: 0,
				RowsAffected: 0,
				Error:        "no such table: baz",
				Time:         0,
			},
		},
	}
	require.Equal(expectedResult, result)
}

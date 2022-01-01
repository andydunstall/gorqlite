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

func TestDataAPIClient_ExecuteThenQueryResults(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	conn := gorqlite.Connect(cluster.Addrs())

	// Create table.
	execResult, err := conn.ExecuteOne(
		"CREATE TABLE foo (id integer not null primary key, name text)",
	)
	require.Nil(err)
	require.Equal("", execResult.Error)

	// Insert one row.
	execResult, err = conn.ExecuteOne(
		`INSERT INTO foo(name) VALUES("fiona")`,
	)
	require.Nil(err)
	require.Equal("", execResult.Error)
	require.Equal(int64(1), execResult.RowsAffected)
	require.Equal(int64(1), execResult.LastInsertId)

	queryResult, err := conn.QueryOne(
		`SELECT name FROM foo WHERE name="fiona"`,
	)
	require.Nil(err)
	require.Equal("", queryResult.Error)
	require.Equal(gorqlite.QueryResult{
		Columns: []string{"name"},
		Values:  [][]interface{}{{"fiona"}},
	}, queryResult)

	// Update one row.
	execResult, err = conn.ExecuteOne(
		`UPDATE foo SET name="justin" WHERE name="fiona"`,
	)
	require.Nil(err)
	require.Equal("", execResult.Error)
	require.Equal(int64(1), execResult.RowsAffected)

	queryResult, err = conn.QueryOne(
		`SELECT name FROM foo WHERE name="justin"`,
	)
	require.Nil(err)
	require.Equal("", queryResult.Error)
	require.Equal(gorqlite.QueryResult{
		Columns: []string{"name"},
		Values:  [][]interface{}{{"justin"}},
	}, queryResult)

	// Delete one row.
	execResult, err = conn.ExecuteOne(
		`DELETE FROM foo WHERE name="justin"`,
	)
	require.Nil(err)
	require.Equal("", execResult.Error)
	require.Equal(int64(1), execResult.RowsAffected)

	queryResult, err = conn.QueryOne(
		`SELECT COUNT(id) AS idCount FROM foo`,
	)
	require.Nil(err)
	require.Equal("", queryResult.Error)
	require.Equal(gorqlite.QueryResult{
		Columns: []string{"idCount"},
		Values:  [][]interface{}{{float64(0)}},
	}, queryResult)

	// Insert multiple rows.
	sql := []string{}
	numRows := 100
	for i := 0; i < numRows; i++ {
		sql = append(sql, fmt.Sprintf(`INSERT INTO foo(name) VALUES("justin-%d")`, i))
	}
	execResults, err := conn.Execute(sql, gorqlite.WithTransaction(true))
	require.Nil(err)
	require.Equal("", execResults.GetFirstError())
	require.Equal(numRows, len(execResults))
	for i := 0; i < numRows; i++ {
		result := execResults[i]
		require.Equal(int64(1), result.RowsAffected)
		require.Equal(int64(i+1), result.LastInsertId)
	}

	queryResult, err = conn.QueryOne(
		`SELECT COUNT(*) AS total FROM foo WHERE name like("justin-%")`,
		gorqlite.WithConsistency("strong"),
	)
	require.Nil(err)
	require.Equal("", queryResult.Error)
	require.Equal(gorqlite.QueryResult{
		Columns: []string{"total"},
		Values:  [][]interface{}{{float64(numRows)}},
	}, queryResult)

	sql = []string{}
	for i := 0; i < numRows; i++ {
		sql = append(sql, fmt.Sprintf(`SELECT name FROM foo WHERE name="justin-%d"`, i))
	}
	queryResults, err := conn.Query(sql)
	require.Nil(err)
	require.Equal("", queryResults.GetFirstError())
	require.Equal(numRows, len(queryResults))
	for i, result := range queryResults {
		require.Equal(gorqlite.QueryResult{
			Columns: []string{"name"},
			Values:  [][]interface{}{{fmt.Sprintf("justin-%d", i)}},
		}, result)
	}

	// TODO
	queryResult, err = conn.QueryOne(
		`SELECT * FROM foo WHERE name like("justin-%")`,
	)
	require.Nil(err)
	for {
		row, ok := queryResult.Next()
		if !ok {
			break
		}
		fmt.Println(row)
	}

	// Drop the table.
	execResult, err = conn.ExecuteOne(`DROP TABLE foo`)
	require.Nil(err)
	require.Equal("", execResult.Error)
	require.Equal(int64(1), execResult.RowsAffected)
}

func TestDataAPIClient_QueryInvalidCommand(t *testing.T) {
	require := require.New(t)

	cluster, err := cluster.OpenCluster(3)
	require.Nil(err)
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	require.True(cluster.WaitForHealthy(ctx))

	conn := gorqlite.Connect(cluster.Addrs())
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

	conn := gorqlite.Connect(cluster.Addrs())
	result, err := conn.Execute([]string{
		"CREATE TABLE foo (id integer not null primary key, bar text)",
		// Table does not exist.
		"INSERT INTO baz (bar) values ('bar')",
	})

	require.Nil(err)

	expectedResult := gorqlite.ExecuteResults{
		{
			LastInsertId: 0,
			RowsAffected: 0,
			Error:        "",
		},
		{
			LastInsertId: 0,
			RowsAffected: 0,
			Error:        "no such table: baz",
		},
	}
	require.Equal(expectedResult, result)
}

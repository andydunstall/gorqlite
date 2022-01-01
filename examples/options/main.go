package main

import (
	"net/http"

	"github.com/dunstall/gorqlite"
	"github.com/dunstall/gorqlite/cluster"
	log "github.com/sirupsen/logrus"
)

func main() {
	cluster, err := cluster.RunDefaultCluster()
	if err != nil {
		log.Fatal(err)
	}
	defer cluster.Close()

	// Open a connection with custom HTTP headers that apply to all requests.
	conn := gorqlite.Connect(cluster.Addrs(), gorqlite.WithHTTPHeaders(http.Header{
		"X-MYHEADER": []string{"my-value"},
	}))

	execResult, err := conn.ExecuteOne(
		"CREATE TABLE foo (id integer not null primary key, name text)",
	)
	if err != nil {
		log.Fatal(err)
	}
	if execResult.Error != "" {
		log.Fatal(execResult.Error)
	}

	// Insert as a transaction.
	execResults, err := conn.Execute([]string{
		`INSERT INTO foo(name) VALUES("foo")`,
		`INSERT INTO foo(name) VALUES("bar")`,
	}, gorqlite.WithTransaction(true))
	if err != nil {
		log.Fatal(err)
	}
	if execResults.HasError() {
		log.Fatal(execResults.GetFirstError())
	}

	// Query with strong consistency.
	sql := []string{
		`SELECT * FROM foo WHERE id="1"`,
		`SELECT * FROM foo WHERE name="bar"`,
	}
	queryResult, err := conn.Query(sql, gorqlite.WithConsistency("strong"))
	if err != nil {
		log.Fatal(err)
	}
	if queryResult.HasError() {
		log.Fatal(queryResult.GetFirstError())
	}
	for _, result := range queryResult {
		for {
			row, ok := result.Next()
			if !ok {
				break
			}

			var id int
			var name string
			if err = row.Scan(&id, &name); err != nil {
				log.Fatal(err)
			}
			log.Info("id:", id)
			log.Info("name:", name)
		}
	}
}

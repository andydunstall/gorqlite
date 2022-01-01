package main

import (
	"fmt"

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

	conn := gorqlite.Open(cluster.Addrs())

	// Create table.
	execResult, err := conn.Execute([]string{
		"CREATE TABLE foo (id integer not null primary key, name text)",
	})
	if err != nil {
		log.Fatal(err)
	}
	if execResult.HasError() {
		log.Fatal(execResult.GetFirstError())
	}

	execResult, err = conn.Execute([]string{
		`INSERT INTO foo(name) VALUES("fiona")`,
	})
	if err != nil {
		log.Fatal(err)
	}
	if execResult.HasError() {
		log.Fatal(execResult.GetFirstError())
	}
	id := execResult[0].LastInsertId
	log.Infof("id of the inserted row: %d", id)

	queryResult, err := conn.Query([]string{
		fmt.Sprintf(`SELECT name FROM foo WHERE id="%d"`, id),
	})
	if err != nil {
		log.Fatal(err)
	}
	if queryResult.HasError() {
		log.Fatal(queryResult.GetFirstError())
	}
	log.Info(queryResult[0])

	execResult, err = conn.Execute([]string{
		`UPDATE foo SET name="justin" WHERE name="fiona"`,
	})
	if err != nil {
		log.Fatal(err)
	}
	if execResult.HasError() {
		log.Fatal(execResult.GetFirstError())
	}
	rowsAffected := execResult[0].RowsAffected
	log.Infof("rows affected: %d", rowsAffected)

	queryResult, err = conn.Query([]string{
		fmt.Sprintf(`SELECT name FROM foo WHERE id="%d"`, id),
	})
	if err != nil {
		log.Fatal(err)
	}
	if queryResult.HasError() {
		log.Fatal(queryResult.GetFirstError())
	}
	log.Info(queryResult[0])
}

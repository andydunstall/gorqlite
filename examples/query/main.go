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

	conn := gorqlite.Connect(cluster.Addrs())

	execResults, err := conn.Execute([]string{
		"CREATE TABLE foo (id integer not null primary key, name text)",
		`INSERT INTO foo(name) VALUES("fiona")`,
	})
	if err != nil {
		log.Fatal(err)
	}
	if execResults.HasError() {
		log.Fatal(execResults.GetFirstError())
	}
	id := execResults[1].LastInsertId
	log.Infof("id of the inserted row: %d", id)

	queryResult, err := conn.QueryOne(
		fmt.Sprintf(`SELECT name FROM foo WHERE id="%d"`, id),
	)
	if err != nil {
		log.Fatal(err)
	}
	if queryResult.Error != "" {
		log.Fatal(queryResult.Error)
	}
	row, _ := queryResult.Next()
	var name string
	if err = row.Scan(&name); err != nil {
		log.Fatal(err)
	}
	log.Info("name:", name)

	execResult, err := conn.ExecuteOne(
		`UPDATE foo SET name="justin" WHERE name="fiona"`,
	)
	if err != nil {
		log.Fatal(err)
	}
	if execResult.Error != "" {
		log.Fatal(execResult.Error)
	}
	rowsAffected := execResult.RowsAffected
	log.Infof("rows affected: %d", rowsAffected)

	queryResult, err = conn.QueryOne(
		fmt.Sprintf(`SELECT name FROM foo WHERE id="%d"`, id),
	)
	if err != nil {
		log.Fatal(err)
	}
	if queryResult.Error != "" {
		log.Fatal(queryResult.Error)
	}
	row, _ = queryResult.Next()
	if err = row.Scan(&name); err != nil {
		log.Fatal(err)
	}
	log.Info("name:", name)
}

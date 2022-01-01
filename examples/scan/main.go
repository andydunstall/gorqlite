package main

import (
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

	execResult, err := conn.Execute([]string{
		"CREATE TABLE foo (id integer not null primary key, name text)",
		`INSERT INTO foo(name) VALUES("foo")`,
		`INSERT INTO foo(name) VALUES("bar")`,
		`INSERT INTO foo(name) VALUES("baz")`,
	})
	if err != nil {
		log.Fatal(err)
	}
	if execResult.HasError() {
		log.Fatal(execResult.GetFirstError())
	}

	queryResult, err := conn.QueryOne("SELECT * FROM foo")
	if err != nil {
		log.Fatal(err)
	}
	for {
		row, ok := queryResult.Next()
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

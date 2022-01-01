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

	execResult, err := conn.ExecuteOne(
		"CREATE TABLE foo (id integer not null primary key, name text)",
	)
	if err != nil {
		log.Fatal(err)
	}
	if execResult.Error != "" {
		log.Fatal(execResult.Error)
	}
}

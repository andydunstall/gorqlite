package gorqlite

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

	sql := []string{
		`INSERT INTO foo(name) VALUES("fiona")`,
		`INSERT INTO bar(name) VALUES("test")`,
	}
	execResult, err = conn.Execute(sql)
	if err != nil {
		log.Fatal(err)
	}
	if execResult.HasError() {
		log.Fatal(execResult.GetFirstError())
	}
	log.Infof("id for first insert: %d", execResult[0].LastInsertId)
	log.Infof("id for second insert: %d", execResult[1].LastInsertId)

	sql = []string{
		`SELECT name FROM foo WHERE id="1"`,
		`SELECT id FROM bar WHERE name="test"`,
	}
	queryResult, err := conn.Query(sql, gorqlite.WithConsistency("strong"))
	if err != nil {
		log.Fatal(err)
	}
	if queryResult.HasError() {
		log.Fatal(queryResult.GetFirstError())
	}
	log.Info(queryResult[0])
	log.Info(queryResult[1])
}

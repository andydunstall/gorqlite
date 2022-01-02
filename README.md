# gorqlite
![build status](https://app.travis-ci.com/dunstall/gorqlite.svg?branch=main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dunstall/gorqlite)](https://pkg.go.dev/github.com/dunstall/gorqlite?tab=doc)

A client library for [rqlite](https://github.com/rqlite/rqlite), the
lightweight, distributed database built on SQLite. This package is designed to
provide a Go interface for the RQLite API endpoints.

## In Progress
See `TODO.md`.

## Examples
### Connect
Connects to an rqlite cluster and creates a table.
```go
addrs := []string{"node-1:8423", "node-2:2841", "node-3"}
conn := gorqlite.Open(addrs)

// Create a table with a single statement.
execResult, err := conn.ExecuteOne(
  "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT, age INTEGER)",
)
if err != nil {
  panic(err)
}
if execResult.Error != "" {
  panic(execResult.Error)
}
```

### Execute Then Query
Inserts a row then selects the row.
```go
// Insert multiple entries in one call.
execResults, err := conn.Execute([]string{
  `INSERT INTO foo(name, age) VALUES(\"fiona\", 20)`,
  `INSERT INTO foo(name, age) VALUES(\"sinead\", 24)`,
})
if err != nil {
  panic(err)
}
if execResults.HasError() {
  panic(execResults.GetFirstError())
}
for _, r := range execResults {
  fmt.Println("id of the inserted row:", r.LastInsertId)
  fmt.Println("rows affected:", r.RowsAffected)
}

// Query the results.
queryResult, err := conn.QueryOne("SELECT * FROM foo")
if err != nil {
  panic(err)
}
if queryResult.Error != "" {
  panic(queryResult.Error)
}

// Scan the results into variables.
for {
  row, ok := queryResult.Next()
  if !ok {
    break
  }

  var id int
  var name string
  if err = row.Scan(&id, &name); err != nil {
    panic(err)
  }
  fmt.Println("ID:", id, "Name:", name)
}
```

### Custom Options
Add default and method override options.
```go
conn := gorqlite.Open(cluster.Addrs(), gorqlite.WithActiveHostRoundRobin(false))

execResult, err := conn.ExecuteOne(
  "CREATE TABLE foo (id integer not null primary key, name text)",
)
if err != nil {
  log.Fatal(err)
}
if execResult.Error != "" {
  log.Fatal(execResult.Error)
}

// Execute the statements within a transaction.
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

// Query the table with strong consistency.
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
```

## Testing
Tests are split into unit tests and system tests.

### Unit Tests
Unit tests are covered in the packages themselves. These are restricted to running in a single thread, with no sleeps and no external accesses (network, files, ...). These can be run quickly with
```go
make test
```

### System Tests
System tests run `gorqlite` against real `rqlite` nodes. These nodes are spun up and down as a cluster within the test itself (see `tests/cluster`). A new cluster is created per test to avoid any side affects (such as some tests for verifying node failover will terminate random nodes in the cluster).

These are disabled by default using the `system` build tag. For additional logging the environment variable `DEBUG=true` can also be used (along with `-v` flag to see the log output).
```go
[DEBUG=true] go test ./... -tags system [-v]
# or just `make system-test`
```

For convenience a docker environment is provided that already has Go and rqlite
installed.
```go
make env
```

The logs for each node can be found in `tests/log`.

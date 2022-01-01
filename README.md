# gorqlite
![build status](https://app.travis-ci.com/dunstall/gorqlite.svg?branch=main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dunstall/gorqlite)](https://pkg.go.dev/github.com/dunstall/gorqlite?tab=doc)

A client library for [rqlite](https://github.com/rqlite/rqlite), the
lightweight, distributed database built on SQLite. This package is designed to
provide a Go interface for the RQLite API endpoints.

The design of this library is based on [rqlite-js](https://github.com/rqlite/rqlite-js).

## In Progress
See `TODO.md`.

## Examples
Contains the same examples as [rqlite-js](https://github.com/rqlite/rqlite-js)
ported into `gorqlite`. See `examples/`, where each example will spin up its
own local cluster (see `Dockerfile` and `make env` for an environment to run
in).

### CREATE TABLE
Connects to an rqlite cluster and creates a table. See `examples/create_table`.
```go
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
```

### Multiple QUERY
Inserts a row then selects the row. See `examples/query`.
```go
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
```

### Custom Options
Add default and method override options. See `examples/options`.
```go
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

# gorqlite
![build status](https://app.travis-ci.com/dunstall/gorqlite.svg?branch=main)

A client library for [rqlite](https://github.com/rqlite/rqlite), the
lightweight, distributed database built on SQLite. This package is designed to
provide a Go interface for the RQLite API endpoints.

The design of this library is based on [rqlite-js](https://github.com/rqlite/rqlite-js).

## In Progress
- [ ] Add travis CI tests to test:
  * `make generate` makes no changes
  * `make fmt` makes no changes
  * `make lint` passes 
  * `make test` passes
  * `make system-test` passes
- [ ] Add API docs (and code level comments)
- [ ] Add support for follow redirects and cache leader (see https://github.com/rqlite/rqlite/blob/master/DOC/DATA_API.md#disabling-request-forwarding)
- [ ] Add queries to result types (such as `QueryResult.Get("name")`)
- [ ] Add system tests for:
  * consistency and transactions
  * node/leader failover/retries
- [ ] Add `/nodes` and `/ready` APIs (see https://github.com/rqlite/rqlite/blob/master/DOC/DIAGNOSTICS.md)
- [ ] Add backup APIs (see https://github.com/rqlite/rqlite/blob/master/DOC/BACKUPS.md)
- [ ] Add long running test with random queries to check for leaks (see go.dev/doc/diagnostic)
- [ ] Review rqlite/rqlite-js, rqlite/gorqlite and rqlite/pyrqlite SDKs for missing tests, invalid handling of requests/responses, etc

## Examples
Contains the same examples as [rqlite-js](https://github.com/rqlite/rqlite-js)
ported into `gorqlite`. See `examples/`, where each example will spin up its
own local cluster (see `Dockerfile` and `make env` for an environment to run
in).

### CREATE TABLE
Connects to an rqlite cluster and creates a table. See `examples/create_table`.
```go
conn := gorqlite.Open(clusterAddrs)

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
```

### Multiple QUERY
Inserts a row then selects the row. See `examples/query`.
```go
execResult, err = conn.Execute([]string{
  `INSERT INTO foo(name) VALUES("fiona")`,
})
if err != nil {
  log.Fatal(err)
}
if execResult.HasError() {
  log.Fatal(execResult.GetFirstError())
}
id := execResult.Results[0].LastInsertId
log.Infof("id of the inserted row: %d" ,id)

queryResult, err := conn.Query([]string{
  fmt.Sprintf(`SELECT name FROM foo WHERE id="%d"`, id),
})
if err != nil {
  log.Fatal(err)
}
if queryResult.HasError() {
  log.Fatal(queryResult.GetFirstError())
}
log.Info(queryResult.Results[0])

execResult, err = conn.Execute([]string{
  `UPDATE foo SET name="justin" WHERE name="fiona"`,
})
if err != nil {
  log.Fatal(err)
}
if execResult.HasError() {
  log.Fatal(execResult.GetFirstError())
}
rowsAffected := execResult.Results[0].RowsAffected
log.Infof("rows affected: %d" ,rowsAffected)

queryResult, err = conn.Query([]string{
  fmt.Sprintf(`SELECT name FROM foo WHERE id="%d"`, id),
})
if err != nil {
  log.Fatal(err)
}
if queryResult.HasError() {
  log.Fatal(queryResult.GetFirstError())
}
log.Info(queryResult.Results[0])
```

### Transactions
Runs multiple queries within a transaction. See `examples/transaction`.
```go
sql := []string{
  `INSERT INTO foo(name) VALUES("fiona")`,
  `INSERT INTO bar(name) VALUES("test")`,
}
execResult, err = conn.Execute(sql, gorqlite.WithTransaction(true))
if err != nil {
  log.Fatal(err)
}
if execResult.HasError() {
  log.Fatal(execResult.GetFirstError())
}
log.Infof("id for first insert: %d", execResult.Results[0].LastInsertId)
log.Infof("id for second insert: %d", execResult.Results[1].LastInsertId)
```

### Consistency
Runs multiple select queries with [strong consistency](https://github.com/rqlite/rqlite/blob/master/DOC/CONSISTENCY.md).
See `examples/consistency`.
```go
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
log.Info(queryResult.Results[0])
log.Info(queryResult.Results[1])
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

Currently this depends on `toxiproxy-server` running (which will be removed soon and integrated into the test itself), used to add network faults in the cluster. A docker environment exists that can be used instead with:
```go
make env
```

The logs for each node can be found in `tests/log`.

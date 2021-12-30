# gorqlite
A client library for [rqlite](https://github.com/rqlite/rqlite) based on [rqlite-js](https://github.com/rqlite/rqlite-js).

## TODO
- [ ] Run go.dev/doc/diagnostic
  * Also setup a long running test that just makes random queries for a few
hours to check for leaks or rare errors
- [ ] Add QueryResult and ExecuteResult lookup methods from rqlite-js
- [ ] Add HTTP fetch error handling from rqlite-js
  * Fail over
  * Retries
  * Redirect
- [ ] Add docs (both MD docs and in the code itself)
- [ ] Check rqlite/gorqlite (and other SDKs) for any missing functionality

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

#### Missing Tests
- [ ] Failover to another node when connected node fails
- [ ] Failover to another node when leader node is partitioned (so still
reachable though cannot handle requests)

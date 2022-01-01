# TODO

## v0.1.0

### APIs
- [ ] Add `/nodes` and `/ready` APIs (see https://github.com/rqlite/rqlite/blob/master/DOC/DIAGNOSTICS.md)
- [ ] Add backup APIs (see https://github.com/rqlite/rqlite/blob/master/DOC/BACKUPS.md)
- [ ] Review rqlite/rqlite-js, rqlite/gorqlite and rqlite/pyrqlite SDKs for missing tests, invalid handling of requests/responses, etc
- [ ] Add better result types (such as `QueryResult.Get("name")`). See `rqlite/gorqlite:QueryResult`.
  * Replace ExecuteResponse with []ExecuteResult and return error if ExecuteResponse.Error != ""
  * Replace QueryResponse with []QueryRows and return error if QueryResponse.Error != ""
  * Add ExecuteOne and QueryOne
- [ ] Add `Leader()` and `Peers()` to API (see `rqlite/gorqlite`)
- [ ] Improve errors
  * If failed to query all nodes, add the error for each of them
- [ ] Maybe add method specific options (such as WithConsistency doesnt appy to status)
  * `Option`
  * `QueryOption`
  * `ExecuteOption`
  * `StatusOption`

### CI
- [ ] Add `make system-test` to CI (which requires installing `rqlite`)

### Docs
- [ ] Add API docs
- [ ] Add better code level comments
- [ ] Add go reference docs (see https://github.com/go-redis/redis for a good example)

### Fault Tolerance
- [ ] Add support for follow redirects and cache leader (see https://github.com/rqlite/rqlite/blob/master/DOC/DATA_API.md#disabling-request-forwarding)
- [ ] Add system tests for
  * consistency and transactions
  * node/leader failover/retries
- [ ] Add support for caching the list of nodes from `/nodes` API and try all of these (such that the user only needs to provide the address of a single node)
- [ ] Add long running test with random queries to check for leaks (see go.dev/doc/diagnostic)

### HTTPS
* [ ] Add HTTPS support

## Future
* Parameterized queries
* `rqlite/system_test` has some useful tests

# TODO

## CI
- [ ] Add `make system-test` to CI (which requires installing `rqlite`)

## Docs
- [ ] Add API docs
- [ ] Add better code level comments

## Fault Tolerance
- [ ] Add support for follow redirects and cache leader (see https://github.com/rqlite/rqlite/blob/master/DOC/DATA_API.md#disabling-request-forwarding)
- [ ] Add system tests for
  * consistency and transactions
  * node/leader failover/retries
- [ ] Add support for caching the list of nodes from `/nodes` API and try all of these (such that the user only needs to provide the address of a single node)
- [ ] Add long running test with random queries to check for leaks (see go.dev/doc/diagnostic)

## HTTPS
* [ ] Add HTTPS support

## Queries
- [ ] Add better result types (such as `QueryResult.Get("name")`)

## APIs
- [ ] Add `/nodes` and `/ready` APIs (see https://github.com/rqlite/rqlite/blob/master/DOC/DIAGNOSTICS.md)
- [ ] Add backup APIs (see https://github.com/rqlite/rqlite/blob/master/DOC/BACKUPS.md)
- [ ] Review rqlite/rqlite-js, rqlite/gorqlite and rqlite/pyrqlite SDKs for missing tests, invalid handling of requests/responses, etc

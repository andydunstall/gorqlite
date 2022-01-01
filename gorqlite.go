package gorqlite

import (
	"context"
	"encoding/json"
)

// Gorqlite is a client for the rqlite API endpoints.
type Gorqlite struct {
	apiClient APIClient
}

// Opens a connection to rqlite.
//
// hosts is a list of addresses (in format host[:port]) for the known
// nodes in the cluster.
//
// opts is a list of default options used for each request (which can be
// overridden on a per request basis by passing opts to each method).
func Open(hosts []string, opts ...Option) *Gorqlite {
	apiClient := newHTTPAPIClient(hosts, opts...)
	return &Gorqlite{
		apiClient,
	}
}

// OpenWithClient opens a connection to rqlite using a custom API client.
func OpenWithClient(apiClient APIClient, opts ...Option) *Gorqlite {
	return &Gorqlite{
		apiClient,
	}
}

// Query runs the given query `sql` statements to rqlite and returns the
// results.
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/DATA_API.md#querying-data.
func (g *Gorqlite) Query(sql []string, opts ...Option) (QueryResponse, error) {
	return g.QueryWithContext(context.Background(), sql, opts...)
}

func (g *Gorqlite) QueryWithContext(ctx context.Context, sql []string, opts ...Option) (QueryResponse, error) {
	body, err := json.Marshal(sql)
	if err != nil {
		return QueryResponse{}, wrapError(err, "query failed: failed to marshal query")
	}
	resp, err := g.apiClient.PostWithContext(ctx, "/db/query", body, opts...)
	if err != nil {
		return QueryResponse{}, wrapError(err, "query failed: request failed")
	}
	defer resp.Body.Close()

	if !isStatusOK(resp.StatusCode) {
		return QueryResponse{}, newError("query failed: invalid status code: %d", resp.StatusCode)
	}

	var results QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return QueryResponse{}, wrapError(err, "query failed: invalid response")
	}

	return results, nil
}

// Execute writes the given `sql` statements to rqlite and returns the execute
// response.
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/DATA_API.md#writing-data.
//
// To enable transactions use `WithTransaction(true)` option.
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/DATA_API.md#transactions.
func (g *Gorqlite) Execute(sql []string, opts ...Option) (ExecuteResponse, error) {
	return g.ExecuteWithContext(context.Background(), sql, opts...)
}

func (g *Gorqlite) ExecuteWithContext(ctx context.Context, sql []string, opts ...Option) (ExecuteResponse, error) {
	body, err := json.Marshal(sql)
	if err != nil {
		return ExecuteResponse{}, wrapError(err, "execute failed: failed to marshal query")
	}
	resp, err := g.apiClient.PostWithContext(ctx, "/db/execute", body, opts...)
	if err != nil {
		return ExecuteResponse{}, wrapError(err, "execute failed: request failed")
	}
	defer resp.Body.Close()

	if !isStatusOK(resp.StatusCode) {
		return ExecuteResponse{}, newError("execute failed: invalid status code: %d", resp.StatusCode)
	}

	var results ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return ExecuteResponse{}, wrapError(err, "execute failed: invalid response")
	}

	return results, nil
}

// Status queries the rqlite status API.
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/DIAGNOSTICS.md#status-and-diagnostics-api.
func (g *Gorqlite) Status(opts ...Option) (Status, error) {
	return g.StatusWithContext(context.Background(), opts...)
}

func (g *Gorqlite) StatusWithContext(ctx context.Context, opts ...Option) (Status, error) {
	resp, err := g.apiClient.GetWithContext(ctx, "/status", opts...)
	if err != nil {
		return Status{}, wrapError(err, "failed to fetch status")
	}
	defer resp.Body.Close()

	if !isStatusOK(resp.StatusCode) {
		return Status{}, newError("failed to fetch status: invalid status code: %d", resp.StatusCode)
	}

	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return Status{}, wrapError(err, "failed to fetch status: invalid response")
	}
	return status, nil
}

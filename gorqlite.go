package gorqlite

import (
	"context"
)

type Gorqlite struct {
	dataClient   *dataAPIClient
	statusClient *statusAPIClient
}

func Open(hosts []string, opts ...Option) *Gorqlite {
	httpClient := newHTTPAPIClient(hosts, opts...)
	dataClient := newDataAPIClient(httpClient)
	statusClient := newStatusAPIClient(httpClient)
	return &Gorqlite{
		dataClient,
		statusClient,
	}
}

func (api *Gorqlite) Query(sql []string) (QueryResponse, error) {
	return api.dataClient.Query(sql)
}

func (api *Gorqlite) QueryWithContext(ctx context.Context, sql []string) (QueryResponse, error) {
	return api.dataClient.QueryWithContext(ctx, sql)
}

func (api *Gorqlite) Execute(sql []string) (ExecuteResponse, error) {
	return api.dataClient.Execute(sql)
}

func (api *Gorqlite) ExecuteWithContext(ctx context.Context, sql []string) (ExecuteResponse, error) {
	return api.dataClient.ExecuteWithContext(ctx, sql)
}

func (api *Gorqlite) Status() (Status, error) {
	return api.statusClient.Status()
}

func (api *Gorqlite) StatusWithContext(ctx context.Context) (Status, error) {
	return api.statusClient.StatusWithContext(ctx)
}

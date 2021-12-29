package gorqlite

import (
	"context"
	"encoding/json"
)

type DataAPIClient struct {
	client APIClient
}

func NewDataAPIClient(addr string) *DataAPIClient {
	return &DataAPIClient{
		client: NewHTTPAPIClient(addr),
	}
}

func NewDataAPIClientWithClient(client APIClient) *DataAPIClient {
	return &DataAPIClient{
		client: client,
	}
}

func (api *DataAPIClient) Query(sql []string) error {
	return api.QueryWithContext(context.Background(), sql)
}

func (api *DataAPIClient) QueryWithContext(ctx context.Context, sql []string) error {
	body, err := json.Marshal(sql)
	if err != nil {
		return WrapError(err, "query failed: failed to marshal query")
	}
	resp, err := api.client.PostWithContext(ctx, "/db/query", body)
	if err != nil {
		return WrapError(err, "query failed: request failed")
	}
	defer resp.Body.Close()

	if !IsStatusOK(resp.StatusCode) {
		return NewError("query failed: invalid status code: %d", resp.StatusCode)
	}

	return nil
}

func (api *DataAPIClient) Execute(sql []string) error {
	return api.ExecuteWithContext(context.Background(), sql)
}

func (api *DataAPIClient) ExecuteWithContext(ctx context.Context, sql []string) error {
	body, err := json.Marshal(sql)
	if err != nil {
		return WrapError(err, "execute failed: failed to marshal query")
	}
	resp, err := api.client.PostWithContext(ctx, "/db/execute", body)
	if err != nil {
		return WrapError(err, "execute failed: request failed")
	}
	defer resp.Body.Close()

	if !IsStatusOK(resp.StatusCode) {
		return NewError("execute failed: invalid status code: %d", resp.StatusCode)
	}

	return nil
}

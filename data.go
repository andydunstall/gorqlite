package gorqlite

import (
	"context"
	"encoding/json"
)

type QueryRows struct {
	Columns []string        `json:"columns,omitempty"`
	Types   []string        `json:"types,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
	Error   string          `json:"error,omitempty"`
	Time    float64         `json:"time,omitempty"`
}

type QueryResponse struct {
	Results []QueryRows `json:"results,omitempty"`
	Error   string      `json:"error,omitempty"`
	Time    float64     `json:"time,omitempty"`
}

func (r *QueryResponse) GetFirstError() string {
	if r.Error != "" {
		return r.Error
	}
	for _, row := range r.Results {
		if row.Error != "" {
			return row.Error
		}
	}
	return ""
}

type ExecuteResult struct {
	LastInsertId int64   `json:"last_insert_id,omitempty"`
	RowsAffected int64   `json:"rows_affected,omitempty"`
	Error        string  `json:"error,omitempty"`
	Time         float64 `json:"time,omitempty"`
}

type ExecuteResponse struct {
	Results []ExecuteResult `json:"results,omitempty"`
	Error   string          `json:"error,omitempty"`
	Time    float64         `json:"time,omitempty"`
}

func (r *ExecuteResponse) GetFirstError() string {
	if r.Error != "" {
		return r.Error
	}
	for _, result := range r.Results {
		if result.Error != "" {
			return result.Error
		}
	}
	return ""
}

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

func (api *DataAPIClient) Query(sql []string) (QueryResponse, error) {
	return api.QueryWithContext(context.Background(), sql)
}

func (api *DataAPIClient) QueryWithContext(ctx context.Context, sql []string) (QueryResponse, error) {
	body, err := json.Marshal(sql)
	if err != nil {
		return QueryResponse{}, WrapError(err, "query failed: failed to marshal query")
	}
	resp, err := api.client.PostWithContext(ctx, "/db/query", body)
	if err != nil {
		return QueryResponse{}, WrapError(err, "query failed: request failed")
	}
	defer resp.Body.Close()

	if !IsStatusOK(resp.StatusCode) {
		return QueryResponse{}, NewError("query failed: invalid status code: %d", resp.StatusCode)
	}

	var results QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return QueryResponse{}, WrapError(err, "query failed: invalid response")
	}

	return results, nil
}

func (api *DataAPIClient) Execute(sql []string) (ExecuteResponse, error) {
	return api.ExecuteWithContext(context.Background(), sql)
}

func (api *DataAPIClient) ExecuteWithContext(ctx context.Context, sql []string) (ExecuteResponse, error) {
	body, err := json.Marshal(sql)
	if err != nil {
		return ExecuteResponse{}, WrapError(err, "execute failed: failed to marshal query")
	}
	resp, err := api.client.PostWithContext(ctx, "/db/execute", body)
	if err != nil {
		return ExecuteResponse{}, WrapError(err, "execute failed: request failed")
	}
	defer resp.Body.Close()

	if !IsStatusOK(resp.StatusCode) {
		return ExecuteResponse{}, NewError("execute failed: invalid status code: %d", resp.StatusCode)
	}

	var results ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return ExecuteResponse{}, WrapError(err, "execute failed: invalid response")
	}

	return results, nil
}

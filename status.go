package gorqlite

import (
	"encoding/json"
	"fmt"
)

type BuildStatus struct {
	Branch    string `json:"branch,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	Commit    string `json:"commit,omitempty"`
	Compiler  string `json:"compiler,omitempty"`
	Version   string `json:"version,omitempty"`
}

type ClusterStatus struct {
	Addr    string `json:"addr,omitempty"`
	APIAddr string `json:"api_addr,omitempty"`
	HTTPS   string `json:"https,omitempty"`
}

type Status struct {
	Build   BuildStatus   `json:"build,omitempty"`
	Cluster ClusterStatus `json:"cluster,omitempty"`
	// TODO(AD)
}

type StatusAPIClient struct {
	client APIClient
}

func NewStatusAPIClient(addr string) *StatusAPIClient {
	return &StatusAPIClient{
		client: NewHTTPAPIClient(addr),
	}
}

func NewStatusAPIClientWithClient(client APIClient) *StatusAPIClient {
	return &StatusAPIClient{
		client: client,
	}
}

func (api *StatusAPIClient) Status() (Status, error) {
	resp, err := api.client.Get("/status")
	if err != nil {
		return Status{}, fmt.Errorf("failed to fetch status: %s", err)
	}
	defer resp.Body.Close()

	if !IsStatusOK(resp.StatusCode) {
		return Status{}, fmt.Errorf("failed to fetch status: invalid status code: %d", resp.StatusCode)
	}

	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return Status{}, fmt.Errorf("failed to fetch status: invalid response: %s", err)
	}
	return status, nil
}

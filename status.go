package gorqlite

import (
	"encoding/json"
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

type HTTPStatus struct {
	Auth     string `json:"auth,omitempty"`
	BindAddr string `json:"bind_addr,omitempty"`
}

type NodeStatus struct {
	StartTime string `json:"start_time,omitempty"`
	Uptime    string `json:"uptime,omitempty"`
}

type OSStatus struct {
	Executable string `json:"executable,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
	Pid        int    `json:"pid,omitempty"`
	Ppid       int    `json:"ppid,omitempty"`
}

type RuntimeStatus struct {
	GoArch       string `json:"GOARCH,omitempty"`
	GoMaxProcs   int    `json:"GOMAXPROCS,omitempty"`
	GoOS         string `json:"GOOS,omitempty"`
	NumCPU       int    `json:"num_cpu,omitempty"`
	NumGoroutine int    `json:"num_goroutine,omitempty"`
	Version      string `json:"version,omitempty"`
}

type Leader struct {
	Addr   string `json:"addr,omitempty"`
	NodeID string `json:"node_id,omitempty"`
}

type Node struct {
	Addr     string `json:"addr,omitempty"`
	ID       string `json:"id,omitempty"`
	Suffrage string `json:"suffrage,omitempty"`
}

type StoreStatus struct {
	Addr             string `json:"addr,omitempty"`
	ApplyTimeout     string `json:"apply_timeout,omitempty"`
	DBAppliedIndex   int    `json:"db_applied_index,omitempty"`
	Dir              string `json:"dir,omitempty"`
	DirSize          int    `json:"dir_size,omitempty"`
	ElectionTimeout  string `json:"election_timeout,omitempty"`
	FSMIndex         int    `json:"fsm_index,omitempty"`
	HeartbeatTimeout string `json:"heartbeat_timeout,omitempty"`
	Leader           Leader `json:"leader,omitempty"`
	NodeID           string `json:"node_id,omitempty"`
	Nodes            []Node `json:"nodes,omitempty"`
}

type Status struct {
	Build   BuildStatus   `json:"build,omitempty"`
	Cluster ClusterStatus `json:"cluster,omitempty"`
	HTTP    HTTPStatus    `json:"http,omitempty"`
	Node    NodeStatus    `json:"node,omitempty"`
	OS      OSStatus      `json:"os,omitempty"`
	Runtime RuntimeStatus `json:"runtime,omitempty"`
	Store   StoreStatus   `json:"store,omitempty"`
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
		return Status{}, WrapError(err, "failed to fetch status")
	}
	defer resp.Body.Close()

	if !IsStatusOK(resp.StatusCode) {
		return Status{}, NewError("failed to fetch status: invalid status code: %d", resp.StatusCode)
	}

	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return Status{}, WrapError(err, "failed to fetch status: invalid response")
	}
	return status, nil
}

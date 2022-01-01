package gorqlite

type StatusBuild struct {
	Branch    string `json:"branch,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	Commit    string `json:"commit,omitempty"`
	Compiler  string `json:"compiler,omitempty"`
	Version   string `json:"version,omitempty"`
}

type StatusCluster struct {
	Addr    string `json:"addr,omitempty"`
	APIAddr string `json:"api_addr,omitempty"`
	HTTPS   string `json:"https,omitempty"`
}

type StatusHTTP struct {
	Auth     string `json:"auth,omitempty"`
	BindAddr string `json:"bind_addr,omitempty"`
}

type StatusNode struct {
	StartTime string `json:"start_time,omitempty"`
	Uptime    string `json:"uptime,omitempty"`
}

type StatusOS struct {
	Executable string `json:"executable,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
	Pid        int    `json:"pid,omitempty"`
	Ppid       int    `json:"ppid,omitempty"`
}

type StatusRuntime struct {
	GoArch       string `json:"GOARCH,omitempty"`
	GoMaxProcs   int    `json:"GOMAXPROCS,omitempty"`
	GoOS         string `json:"GOOS,omitempty"`
	NumCPU       int    `json:"num_cpu,omitempty"`
	NumGoroutine int    `json:"num_goroutine,omitempty"`
	Version      string `json:"version,omitempty"`
}

type LeaderInfo struct {
	Addr   string `json:"addr,omitempty"`
	NodeID string `json:"node_id,omitempty"`
}

type NodeInfo struct {
	Addr     string `json:"addr,omitempty"`
	ID       string `json:"id,omitempty"`
	Suffrage string `json:"suffrage,omitempty"`
}

type StatusStore struct {
	Addr             string     `json:"addr,omitempty"`
	ApplyTimeout     string     `json:"apply_timeout,omitempty"`
	DBAppliedIndex   int        `json:"db_applied_index,omitempty"`
	Dir              string     `json:"dir,omitempty"`
	DirSize          int        `json:"dir_size,omitempty"`
	ElectionTimeout  string     `json:"election_timeout,omitempty"`
	FSMIndex         int        `json:"fsm_index,omitempty"`
	HeartbeatTimeout string     `json:"heartbeat_timeout,omitempty"`
	Leader           LeaderInfo `json:"leader,omitempty"`
	NodeID           string     `json:"node_id,omitempty"`
	Nodes            []NodeInfo `json:"nodes,omitempty"`
}

type Status struct {
	Build   StatusBuild   `json:"build,omitempty"`
	Cluster StatusCluster `json:"cluster,omitempty"`
	HTTP    StatusHTTP    `json:"http,omitempty"`
	Node    StatusNode    `json:"node,omitempty"`
	OS      StatusOS      `json:"os,omitempty"`
	Runtime StatusRuntime `json:"runtime,omitempty"`
	Store   StatusStore   `json:"store,omitempty"`
}

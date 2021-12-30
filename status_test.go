package gorqlite

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/dunstall/gorqlite/mocks"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	// statusV6_7_0 is a capture of the status API taken from a node ran with
	// `tests/cluster` running rqlite v6.7.0.
	statusV6_7_0 = `
{
    "build": {
        "branch": "master",
        "build_time": "2021-10-22T18:32:15-0400",
        "commit": "eb6da8f22cfd57d9f46cc31179de7d0cefe2f962",
        "compiler": "gc",
        "version": "v6.7.0"
    },
    "cluster": {
        "addr": "0.0.0.0:7001",
        "api_addr": "0.0.0.0:6001",
        "https": "false"
    },
    "http": {
        "auth": "disabled",
        "bind_addr": "[::]:4001",
        "cluster": {
            "local_node_addr": "0.0.0.0:7001",
            "timeout": 30000000000
        }
    },
    "node": {
        "start_time": "2021-12-20T21:05:05.943343316Z",
        "uptime": "15.696393ms"
    },
    "os": {
        "executable": "/usr/local/bin/rqlited",
        "hostname": "12bcec7079d5",
        "page_size": 4096,
        "pid": 4541,
        "ppid": 4523
    },
    "runtime": {
        "GOARCH": "amd64",
        "GOMAXPROCS": 8,
        "GOOS": "linux",
        "num_cpu": 8,
        "num_goroutine": 12,
        "version": "go1.16"
    },
    "store": {
        "addr": "0.0.0.0:7001",
        "apply_timeout": "10s",
        "db_applied_index": 0,
        "db_conf": {
            "fk_constraints": false,
            "memory": true
        },
        "dir": "/tmp/node-13591731195",
        "dir_size": 32768,
        "election_timeout": "1s",
        "fsm_index": 0,
        "heartbeat_timeout": "1s",
        "leader": {
          "addr": "0.0.0.0:7002",
            "node_id": "2"
        },
        "node_id": "1",
        "nodes": [
            {
                "addr": "0.0.0.0:7001",
                "id": "1",
                "suffrage": "Voter"
            },
            {
                "addr": "0.0.0.0:7002",
                "id": "2",
                "suffrage": "Voter"
            },
            {
                "addr": "0.0.0.0:7003",
                "id": "3",
                "suffrage": "Voter"
            }
        ],
        "raft": {
            "applied_index": 0,
            "bolt": {
                "free_alloc": 8192,
                "free_list_inuse": 32,
                "num_free_pages": 0,
                "num_pending_pages": 2,
                "num_tx_open": 0,
                "num_tx_read": 8,
                "tx_stats": {
                    "cursor_count": 32,
                    "node_count": 9,
                    "node_deref": 0,
                    "page_alloc": 40960,
                    "page_count": 10,
                    "rebalance": 0,
                    "rebalance_time": 0,
                    "spill": 5,
                    "spill_time": 32078,
                    "split": 0,
                    "write": 15,
                    "write_time": 3435589
                }
            },
            "commit_index": 0,
            "fsm_pending": 0,
            "last_contact": "never",
            "last_log_index": 1,
            "last_log_term": 1,
            "last_snapshot_index": 0,
            "last_snapshot_term": 0,
            "latest_configuration": "[{Suffrage:Voter ID:1 Address:0.0.0.0:7001}]",
            "latest_configuration_index": 0,
            "log_size": 32768,
            "num_peers": 0,
            "protocol_version": 3,
            "protocol_version_max": 3,
            "protocol_version_min": 0,
            "snapshot_version_max": 1,
            "snapshot_version_min": 0,
            "state": "Follower",
            "term": 1
        },
        "request_marshaler": {
            "compression_batch": 5,
            "compression_size": 150,
            "force_compression": false
        },
        "snapshot_interval": 30000000000,
        "snapshot_threshold": 8192,
        "sqlite3": {
            "compile_options": [
                "COMPILER=gcc-9.3.0",
                "DEFAULT_WAL_SYNCHRONOUS=1",
                "ENABLE_DBSTAT_VTAB",
                "ENABLE_FTS3",
                "ENABLE_FTS3_PARENTHESIS",
                "ENABLE_JSON1",
                "ENABLE_RTREE",
                "ENABLE_UPDATE_DELETE_LIMIT",
                "OMIT_DEPRECATED",
                "OMIT_LOAD_EXTENSION",
                "OMIT_SHARED_CACHE",
                "SYSTEM_MALLOC",
                "THREADSAFE=1"
            ],
            "conn_pool_stats": {
                "ro": {
                    "idle": 1,
                    "in_use": 0,
                    "max_idle_closed": 0,
                    "max_idle_time_closed": 0,
                    "max_lifetime_closed": 0,
                    "max_open_connections": 0,
                    "open_connections": 1,
                    "wait_count": 0,
                    "wait_duration": 0
                },
                "rw": {
                    "idle": 1,
                    "in_use": 0,
                    "max_idle_closed": 0,
                    "max_idle_time_closed": 0,
                    "max_lifetime_closed": 0,
                    "max_open_connections": 1,
                    "open_connections": 1,
                    "wait_count": 0,
                    "wait_duration": 0
                }
            },
            "db_size": 0,
            "mem_stats": {
                "cache_size": -2000,
                "freelist_count": 0,
                "hard_heap_limit": 0,
                "max_page_count": 1073741823,
                "page_count": 0,
                "page_size": 4096,
                "soft_heap_limit": 0
            },
            "path": ":memory:",
            "ro_dsn": "file:/NlDCoACqBdhpEhaCGkoM?mode=ro&vfs=memdb&_txlock=deferred&_fk=false",
            "rw_dsn": "file:/NlDCoACqBdhpEhaCGkoM?mode=rw&vfs=memdb&_txlock=immediate&_fk=false",
            "version": "3.36.0"
        },
        "trailing_logs": 10240
    }
}
  `
)

func TestStatusAPIClient_v6_7_0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(statusV6_7_0)),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().GetWithContext(gomock.Any(), "/status").Return(resp, nil)

	expected := Status{
		Build: BuildStatus{
			Branch:    "master",
			BuildTime: "2021-10-22T18:32:15-0400",
			Commit:    "eb6da8f22cfd57d9f46cc31179de7d0cefe2f962",
			Compiler:  "gc",
			Version:   "v6.7.0",
		},
		Cluster: ClusterStatus{
			Addr:    "0.0.0.0:7001",
			APIAddr: "0.0.0.0:6001",
			HTTPS:   "false",
		},
		HTTP: HTTPStatus{
			Auth:     "disabled",
			BindAddr: "[::]:4001",
		},
		Node: NodeStatus{
			StartTime: "2021-12-20T21:05:05.943343316Z",
			Uptime:    "15.696393ms",
		},
		OS: OSStatus{
			Executable: "/usr/local/bin/rqlited",
			Hostname:   "12bcec7079d5",
			PageSize:   4096,
			Pid:        4541,
			Ppid:       4523,
		},
		Runtime: RuntimeStatus{
			GoArch:       "amd64",
			GoMaxProcs:   8,
			GoOS:         "linux",
			NumCPU:       8,
			NumGoroutine: 12,
			Version:      "go1.16",
		},
		Store: StoreStatus{
			Addr:             "0.0.0.0:7001",
			ApplyTimeout:     "10s",
			DBAppliedIndex:   0,
			Dir:              "/tmp/node-13591731195",
			DirSize:          32768,
			ElectionTimeout:  "1s",
			FSMIndex:         0,
			HeartbeatTimeout: "1s",
			Leader: Leader{
				Addr:   "0.0.0.0:7002",
				NodeID: "2",
			},
			NodeID: "1",
			Nodes: []Node{
				Node{
					Addr:     "0.0.0.0:7001",
					ID:       "1",
					Suffrage: "Voter",
				},
				Node{
					Addr:     "0.0.0.0:7002",
					ID:       "2",
					Suffrage: "Voter",
				},
				Node{
					Addr:     "0.0.0.0:7003",
					ID:       "3",
					Suffrage: "Voter",
				},
			},
		},
	}

	statusClient := NewStatusAPIClient(apiClient)
	status, err := statusClient.Status()
	require.Nil(t, err)
	require.Equal(t, expected, status)
}

func TestStatusAPIClient_BadStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resp := &http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().GetWithContext(gomock.Any(), "/status").Return(resp, nil)

	statusClient := NewStatusAPIClient(apiClient)
	_, err := statusClient.Status()
	require.Error(t, err)
}

func TestStatusAPIClient_NetworkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().GetWithContext(gomock.Any(), "/status").Return(nil, fmt.Errorf("network err"))

	statusClient := NewStatusAPIClient(apiClient)
	_, err := statusClient.Status()
	require.Error(t, err)
}

package cluster

import (
	"io/ioutil"
	"net"
)

type nodeConfig struct {
	APIAddr           string
	APIAdvAddr        string
	RaftAddr          string
	RaftAdvAddr       string
	Dir               string
	RaftVoter         bool
	RaftSnapThreshold int
	RaftSnapInterval  string
	Join              string
}

func defaultNodeConfig() *nodeConfig {
	return &nodeConfig{
		APIAddr:           mustRandomAddr(),
		APIAdvAddr:        mustRandomAddr(),
		RaftAddr:          mustRandomAddr(),
		RaftAdvAddr:       mustRandomAddr(),
		Dir:               mustRandomPath(),
		RaftVoter:         true,
		RaftSnapThreshold: 8192,
		RaftSnapInterval:  "1s",
		Join:              "",
	}
}

// NodeOption overrides the default node configuration.
type NodeOption func(conf *nodeConfig)

func WithAPIAddr(addr string) NodeOption {
	return func(conf *nodeConfig) {
		conf.APIAddr = addr
	}
}

func WithAPIAdvAddr(addr string) NodeOption {
	return func(conf *nodeConfig) {
		conf.APIAdvAddr = addr
	}
}

func WithRaftAddr(addr string) NodeOption {
	return func(conf *nodeConfig) {
		conf.RaftAddr = addr
	}
}

func WithRaftAdvAddr(addr string) NodeOption {
	return func(conf *nodeConfig) {
		conf.RaftAdvAddr = addr
	}
}

func WithDir(dir string) NodeOption {
	return func(conf *nodeConfig) {
		conf.Dir = dir
	}
}

func WithRaftVoter(voter bool) NodeOption {
	return func(conf *nodeConfig) {
		conf.RaftVoter = voter
	}
}

func WithRaftSnapThreshold(threshold int) NodeOption {
	return func(conf *nodeConfig) {
		conf.RaftSnapThreshold = threshold
	}
}

func WithRaftSnapInterval(interval string) NodeOption {
	return func(conf *nodeConfig) {
		conf.RaftSnapInterval = interval
	}
}

// WithJoin requests the node joins an existing cluster at addr.
func WithJoin(addr string) NodeOption {
	return func(conf *nodeConfig) {
		conf.Join = addr
	}
}

func mustRandomAddr() string {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(wrapError(err, "failed to get random port"))
	}
	addr := listener.Addr().String()
	listener.Close()
	return addr
}

func mustRandomPath() string {
	dir, err := ioutil.TempDir("", "rqlite")
	if err != nil {
		panic(wrapError(err, "failed to open tmpdir"))
	}
	return dir
}

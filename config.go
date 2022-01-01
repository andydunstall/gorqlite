package gorqlite

import (
	"net/http"
)

type Config struct {
	// ActiveHostRoundRobin load balances requests among all known nodes in
	// a round robin strategy. Enabled by default but should disable if nodes
	// are configured in preference order (such that it tries the first unless
	// that fails).
	ActiveHostRoundRobin bool

	// HTTPHeaders adds extra headers to each request.
	HTTPHeaders http.Header

	// Transaction sets the `transaction` query parameter when set.
	// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/DATA_API.md#transactions.
	Transaction bool

	// Consistency sets the `level` query parameter if set, which can be one of
	// "none", "weak" or "strong.
	// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/CONSISTENCY.md.
	Consistency string

	// transport is the underlying HTTP transport used for requests. This is
	// only expected to be used for unit tests to mock the transport.
	transport http.RoundTripper

	// clock is the underlying clock used to sleep. This is only expected to be
	// used for unit tests to mock the clock.
	clock clock
}

// DefaultConfig returns the default configuration which is used as a base
// for all `Option` overrides.
func DefaultConfig() *Config {
	return &Config{
		ActiveHostRoundRobin: true,
		HTTPHeaders:          make(http.Header),
		Transaction:          false,
		Consistency:          "",
		transport:            http.DefaultTransport,
		clock:                &systemClock{},
	}
}

// Option updates the configuration. Options can be set in `gorqlite.Open`
// to set the defaults and overriden per request (which will only apply to
// that request).
type Option func(conf *Config)

func WithActiveHostRoundRobin(roundRobin bool) Option {
	return func(conf *Config) {
		conf.ActiveHostRoundRobin = roundRobin
	}
}

func WithHTTPHeaders(headers http.Header) Option {
	return func(conf *Config) {
		conf.HTTPHeaders = headers
	}
}

func WithTransaction(transaction bool) Option {
	return func(conf *Config) {
		conf.Transaction = transaction
	}
}

func WithConsistency(consistency string) Option {
	return func(conf *Config) {
		conf.Consistency = consistency
	}
}

func withTransport(transport http.RoundTripper) Option {
	return func(conf *Config) {
		conf.transport = transport
	}
}

func withClock(clock clock) Option {
	return func(conf *Config) {
		conf.clock = clock
	}
}

package gorqlite

import (
	"net/http"
)

type config struct {
	ActiveHostRoundRobin bool

	HTTPHeaders http.Header

	Transaction bool

	Consistency string

	// transport is the underlying HTTP transport used for requests. This is
	// only expected to be used for unit tests to mock the transport.
	transport http.RoundTripper

	// clock is the underlying clock used to sleep. This is only expected to be
	// used for unit tests to mock the clock.
	clock clock
}

// defaultConfig returns the default configuration which is used as a base
// for all `Option` overrides.
func defaultConfig() *config {
	return &config{
		ActiveHostRoundRobin: true,
		HTTPHeaders:          make(http.Header),
		Transaction:          false,
		Consistency:          "",
		transport:            http.DefaultTransport,
		clock:                &systemClock{},
	}
}

// Options overrides the default configuration used for each request.
//
// A set of default options can be set in gorqlite.Open, and each method
// accepts options to override the defaults for that request only.
type Option func(conf *config)

// WithActiveHostRoundRobin load balances requests among all known nodes in
// a round robin strategy if enabled. Otherwise will always try nodes in order
// until one works.
//
// Enabled by default.
func WithActiveHostRoundRobin(enabled bool) Option {
	return func(conf *config) {
		conf.ActiveHostRoundRobin = enabled
	}
}

// WithHTTPHeaders adds HTTP headers to the request.
func WithHTTPHeaders(headers http.Header) Option {
	return func(conf *config) {
		conf.HTTPHeaders = headers
	}
}

// WithTransaction sets the `transaction` query parameter when enabled.
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/DATA_API.md#transactions.
//
// Disabed by default.
func WithTransaction(transaction bool) Option {
	return func(conf *config) {
		conf.Transaction = transaction
	}
}

// WithConsistency sets the `level` query parameter if set, otherwise it is not
// set (so rqlite will default to weak consistency.
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/CONSISTENCY.md.
func WithConsistency(consistency string) Option {
	return func(conf *config) {
		conf.Consistency = consistency
	}
}

func withTransport(transport http.RoundTripper) Option {
	return func(conf *config) {
		conf.transport = transport
	}
}

func withClock(clock clock) Option {
	return func(conf *config) {
		conf.clock = clock
	}
}

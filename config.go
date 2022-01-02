package gorqlite

type config struct {
	ActiveHostRoundRobin bool
}

// defaultConfig returns the default configuration which is used as a base
// for all `Option` overrides.
func defaultConfig() *config {
	return &config{
		ActiveHostRoundRobin: true,
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

type queryConfig struct {
	Consistency string
}

func defaultQueryConfig() *queryConfig {
	return &queryConfig{
		Consistency: "",
	}
}

type QueryOption func(conf *queryConfig)

// WithConsistency sets the level query parameter if set, otherwise it is not
// set (so rqlite will default to weak consistency).
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/CONSISTENCY.md.
func WithConsistency(consistency string) QueryOption {
	return func(conf *queryConfig) {
		conf.Consistency = consistency
	}
}

type executeConfig struct {
	Transaction bool
}

func defaultExecuteConfig() *executeConfig {
	return &executeConfig{
		Transaction: false,
	}
}

type ExecuteOption func(conf *executeConfig)

// WithTransaction sets the transaction query parameter when enabled.
// See https://github.com/rqlite/rqlite/blob/cc74ab0af7c128582b7f0fd380033d43e642a121/DOC/DATA_API.md#transactions.
//
// Disabed by default.
func WithTransaction(transaction bool) ExecuteOption {
	return func(conf *executeConfig) {
		conf.Transaction = transaction
	}
}

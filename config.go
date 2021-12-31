package gorqlite

import (
	"net/http"
)

type Config struct {
	ActiveHostRoundRobin bool
	HTTPHeaders          http.Header
	Transport            http.RoundTripper
	RedirectAttempts     int
	clock                Clock
}

type Option func(conf *Config)

func DefaultConfig() *Config {
	return &Config{
		ActiveHostRoundRobin: true,
		HTTPHeaders:          make(http.Header),
		Transport:            http.DefaultTransport,
		RedirectAttempts:     10,
		clock:                &SystemClock{},
	}
}

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

func WithTransport(transport http.RoundTripper) Option {
	return func(conf *Config) {
		conf.Transport = transport
	}
}

func WithRedirectAttempts(redirectAttempts int) Option {
	return func(conf *Config) {
		conf.RedirectAttempts = redirectAttempts
	}
}

func WithClock(clock Clock) Option {
	return func(conf *Config) {
		conf.clock = clock
	}
}

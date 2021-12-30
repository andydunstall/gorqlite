//go:generate mockgen -source http_api.go -destination mocks/mock_http_api.go

package gorqlite

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Define here to generate mock.
type RoundTripper interface {
	http.RoundTripper
}

type Clock interface {
	Sleep(d time.Duration)
}

type SystemClock struct{}

func (c *SystemClock) Sleep(d time.Duration) {
	<-time.After(d)
}

type HTTPConfig struct {
	ActiveHostRoundRobin bool
	HTTPHeaders          http.Header
	Transport            http.RoundTripper
	RedirectAttempts     int
	clock                Clock
}

type HTTPOption func(conf *HTTPConfig)

func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		ActiveHostRoundRobin: true,
		HTTPHeaders:          make(http.Header),
		Transport:            http.DefaultTransport,
		RedirectAttempts:     10,
		clock:                &SystemClock{},
	}
}

func WithActiveHostRoundRobin(roundRobin bool) HTTPOption {
	return func(conf *HTTPConfig) {
		conf.ActiveHostRoundRobin = roundRobin
	}
}

func WithHTTPHeaders(headers http.Header) HTTPOption {
	return func(conf *HTTPConfig) {
		conf.HTTPHeaders = headers
	}
}

func WithTransport(transport http.RoundTripper) HTTPOption {
	return func(conf *HTTPConfig) {
		conf.Transport = transport
	}
}

func WithRedirectAttempts(redirectAttempts int) HTTPOption {
	return func(conf *HTTPConfig) {
		conf.RedirectAttempts = redirectAttempts
	}
}

func WithClock(clock Clock) HTTPOption {
	return func(conf *HTTPConfig) {
		conf.clock = clock
	}
}

type HTTPAPIClient struct {
	hosts           []string
	activeHostIndex int
	client          *http.Client
	conf            *HTTPConfig
}

func NewHTTPAPIClient(hosts []string, opts ...HTTPOption) *HTTPAPIClient {
	conf := DefaultHTTPConfig()
	for _, opt := range opts {
		opt(conf)
	}
	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: conf.Transport,
	}
	return &HTTPAPIClient{
		hosts:           hosts,
		activeHostIndex: 0,
		client:          client,
		conf:            conf,
	}
}

func (api *HTTPAPIClient) Get(path string) (*http.Response, error) {
	return api.fetch(context.Background(), http.MethodGet, path, nil)
}

func (api *HTTPAPIClient) GetWithContext(ctx context.Context, path string) (*http.Response, error) {
	return api.fetch(ctx, http.MethodGet, path, nil)
}

func (api *HTTPAPIClient) Post(path string, body []byte) (*http.Response, error) {
	return api.fetch(context.Background(), http.MethodPost, path, body)
}

func (api *HTTPAPIClient) PostWithContext(ctx context.Context, path string, body []byte) (*http.Response, error) {
	return api.fetch(ctx, http.MethodPost, path, body)
}

func (api *HTTPAPIClient) fetch(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	defer api.rotateActiveHost()

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	redirectAttempts := 0
	retryAttempts := 0
	u := &url.URL{
		Scheme: "http",
		// Host set per retry.
		Host: "",
		Path: path,
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, WrapError(err, "failed to fetch: invalid request")
	}
	if api.conf.HTTPHeaders != nil {
		req.Header = api.conf.HTTPHeaders
	}

	for {
		// TODO(AD) Add a use leader flag. Otherwise will be redirected when leader
		// is known.
		activeHost := api.activeHost()
		if activeHost == "" {
			return nil, NewError("failed to fetch: no addresses given")
		}
		req.URL.Host = activeHost
		req.Host = activeHost

		resp, err := api.client.Do(req)
		if err != nil || isRetryable(resp.StatusCode) {
			if retryAttempts >= (len(api.hosts) * 3) {
				if err != nil {
					return nil, WrapError(err, "failed to fetch: max retries exceeded")
				}
				return nil, NewError("failed to fetch: max retries exceeded: status: %d", resp.StatusCode)
			}

			api.conf.clock.Sleep(waitTimeExponential(retryAttempts, time.Millisecond*100))

			api.rotateActiveHost()
			retryAttempts++
			continue
		}

		if !isRedirect(resp.StatusCode) {
			return resp, nil
		}

		if redirectAttempts >= api.conf.RedirectAttempts {
			return nil, NewError("failed to fetch: max redirects exceeded (%d)", api.conf.RedirectAttempts)
		}
		u, err := url.Parse(resp.Header.Get("location"))
		if err != nil {
			return nil, WrapError(err, "failed to fetch: invalid redirect url")
		}
		req.URL = u
		redirectAttempts++

		// TODO(AD) If redirected store new leader?.
	}
}

func (api *HTTPAPIClient) activeHost() string {
	if 0 <= api.activeHostIndex && api.activeHostIndex < len(api.hosts) {
		return api.hosts[api.activeHostIndex]
	}
	return ""
}

func (api *HTTPAPIClient) rotateActiveHost() {
	if api.conf.ActiveHostRoundRobin {
		api.activeHostIndex = ((api.activeHostIndex + 1) % len(api.hosts))
	}
}

func isRedirect(statusCode int) bool {
	redirectCodes := []int{http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther}
	for _, redirectCode := range redirectCodes {
		if statusCode == redirectCode {
			return true
		}
	}
	return false
}

func isRetryable(statusCode int) bool {
	retryableCodes := []int{
		http.StatusRequestTimeout,
		http.StatusRequestEntityTooLarge,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}
	for _, retryableCode := range retryableCodes {
		if statusCode == retryableCode {
			return true
		}
	}
	return false

}

func waitTimeExponential(attempt int, base time.Duration) time.Duration {
	// 2^attempt * base
	return time.Duration(1<<attempt) * base
}

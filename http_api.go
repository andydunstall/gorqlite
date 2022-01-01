//go:generate mockgen -source http_api.go -destination mocks/http_api/mock_http_api.go

package gorqlite

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Redefine here to generate mock.
type roundTripper interface {
	http.RoundTripper
}

type clock interface {
	Sleep(d time.Duration)
}

type systemClock struct{}

func (c *systemClock) Sleep(d time.Duration) {
	<-time.After(d)
}

type httpAPIClient struct {
	hosts           []string
	activeHostIndex int
	client          *http.Client
	conf            *config
}

func newHTTPAPIClient(hosts []string, opts ...Option) *httpAPIClient {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	client := &http.Client{
		Transport: conf.transport,
	}
	return &httpAPIClient{
		hosts:           hosts,
		activeHostIndex: 0,
		client:          client,
		conf:            conf,
	}
}

func (api *httpAPIClient) Get(path string, opts ...Option) (*http.Response, error) {
	return api.fetch(context.Background(), http.MethodGet, path, nil, opts...)
}

func (api *httpAPIClient) GetWithContext(ctx context.Context, path string, opts ...Option) (*http.Response, error) {
	return api.fetch(ctx, http.MethodGet, path, nil, opts...)
}

func (api *httpAPIClient) Post(path string, body []byte, opts ...Option) (*http.Response, error) {
	return api.fetch(context.Background(), http.MethodPost, path, body, opts...)
}

func (api *httpAPIClient) PostWithContext(ctx context.Context, path string, body []byte, opts ...Option) (*http.Response, error) {
	return api.fetch(ctx, http.MethodPost, path, body, opts...)
}

func (api *httpAPIClient) fetch(ctx context.Context, method, path string, body []byte, opts ...Option) (*http.Response, error) {
	defer api.rotateActiveHost(false)

	// Apply overrides to a copy of api conf.
	conf := *api.conf
	for _, opt := range opts {
		opt(&conf)
	}

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	retryAttempts := 0
	query := url.Values{}
	if conf.Transaction {
		query.Add("transaction", "")
	}
	if conf.Consistency != "" {
		query.Add("level", conf.Consistency)
	}
	u := &url.URL{
		Scheme: "http",
		// Host set per retry.
		Host:     "",
		Path:     path,
		RawQuery: query.Encode(),
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, wrapError(err, "failed to fetch: invalid request")
	}
	if conf.HTTPHeaders != nil {
		req.Header = conf.HTTPHeaders
	}

	for {
		activeHost := api.activeHost()
		if activeHost == "" {
			return nil, newError("failed to fetch: no addresses given")
		}
		req.URL.Host = activeHost
		req.Host = activeHost

		resp, err := api.client.Do(req)
		if err == nil && isStatusOK(resp.StatusCode) {
			return resp, nil
		}

		if err == nil && !isRetryable(resp.StatusCode) {
			return nil, newError("failed to fetch: bad status code: status: %d", resp.StatusCode)
		}

		if retryAttempts >= (len(api.hosts) * 3) {
			if err != nil {
				return nil, wrapError(err, "failed to fetch: max retries exceeded")
			}
			return nil, newError("failed to fetch: max retries exceeded: status: %d", resp.StatusCode)
		}

		conf.clock.Sleep(waitTimeExponential(retryAttempts, time.Millisecond*100))

		// Force rotate even if round robin is disabled.
		api.rotateActiveHost(true)
		retryAttempts++
	}
}

func (api *httpAPIClient) activeHost() string {
	if 0 <= api.activeHostIndex && api.activeHostIndex < len(api.hosts) {
		return api.hosts[api.activeHostIndex]
	}
	return ""
}

func (api *httpAPIClient) rotateActiveHost(force bool) {
	if api.conf.ActiveHostRoundRobin || force {
		api.activeHostIndex = ((api.activeHostIndex + 1) % len(api.hosts))
	}
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

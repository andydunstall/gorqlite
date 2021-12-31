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

type httpAPIClient struct {
	hosts           []string
	activeHostIndex int
	client          *http.Client
	conf            *Config
}

func newHTTPAPIClient(hosts []string, opts ...Option) *httpAPIClient {
	conf := DefaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: conf.Transport,
	}
	return &httpAPIClient{
		hosts:           hosts,
		activeHostIndex: 0,
		client:          client,
		conf:            conf,
	}
}

func (api *httpAPIClient) Get(path string) (*http.Response, error) {
	return api.fetch(context.Background(), http.MethodGet, path, nil)
}

func (api *httpAPIClient) GetWithContext(ctx context.Context, path string) (*http.Response, error) {
	return api.fetch(ctx, http.MethodGet, path, nil)
}

func (api *httpAPIClient) Post(path string, body []byte) (*http.Response, error) {
	return api.fetch(context.Background(), http.MethodPost, path, body)
}

func (api *httpAPIClient) PostWithContext(ctx context.Context, path string, body []byte) (*http.Response, error) {
	return api.fetch(ctx, http.MethodPost, path, body)
}

func (api *httpAPIClient) fetch(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
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

func (api *httpAPIClient) activeHost() string {
	if 0 <= api.activeHostIndex && api.activeHostIndex < len(api.hosts) {
		return api.hosts[api.activeHostIndex]
	}
	return ""
}

func (api *httpAPIClient) rotateActiveHost() {
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

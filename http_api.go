//go:generate mockgen -source http_api.go -destination mocks/mock_http_api.go

package gorqlite

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

// Define here to generate mock.
type RoundTripper interface {
	http.RoundTripper
}

type HTTPConfig struct {
	HTTPHeaders      http.Header
	Transport        http.RoundTripper
	RedirectAttempts int
}

type HTTPOption func(conf *HTTPConfig)

func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		HTTPHeaders:      make(http.Header),
		Transport:        http.DefaultTransport,
		RedirectAttempts: 10,
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

type HTTPAPIClient struct {
	addr   string
	client *http.Client
	conf   *HTTPConfig
}

func NewHTTPAPIClient(addr string, opts ...HTTPOption) *HTTPAPIClient {
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
		addr:   addr,
		client: client,
		conf:   conf,
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
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	redirectAttempts := 0
	u := &url.URL{
		Scheme: "http",
		Host:   api.addr,
		Path:   path,
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, WrapError(err, "failed to fetch: invalid request")
	}
	if api.conf.HTTPHeaders != nil {
		req.Header = api.conf.HTTPHeaders
	}

	for {
		resp, err := api.client.Do(req)
		if err != nil {
			return nil, WrapError(err, "failed to fetch")
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

func isRedirect(statusCode int) bool {
	redirectCodes := []int{http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther}
	for _, redirectCode := range redirectCodes {
		if statusCode == redirectCode {
			return true
		}
	}
	return false
}

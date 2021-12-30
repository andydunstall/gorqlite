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
	HTTPHeaders http.Header
	Transport   http.RoundTripper
}

type HTTPOption func(conf *HTTPConfig)

func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		HTTPHeaders: make(http.Header),
		Transport:   http.DefaultTransport,
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

	url := &url.URL{
		Scheme: "http",
		Host:   api.addr,
		Path:   path,
	}
	req, err := http.NewRequestWithContext(ctx, method, url.String(), reqBody)
	if err != nil {
		return nil, WrapError(err, "failed to fetch: invalid request")
	}
	if api.conf.HTTPHeaders != nil {
		req.Header = api.conf.HTTPHeaders
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, WrapError(err, "failed to fetch")
	}
	return resp, nil
}

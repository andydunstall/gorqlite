//go:generate mockgen -source http_api.go -destination mocks/mock_http_api.go

package gorqlite

import (
	"net/http"
	"net/url"
)

type HTTPConfig struct {
	HTTPHeaders http.Header
}

type HTTPOption func(conf *HTTPConfig)

func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		HTTPHeaders: make(http.Header),
	}
}

func WithHTTPHeaders(headers http.Header) HTTPOption {
	return func(conf *HTTPConfig) {
		conf.HTTPHeaders = headers
	}
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPAPIClient struct {
	addr   string
	client httpClient
	conf   *HTTPConfig
}

func NewHTTPAPIClient(addr string, opts ...HTTPOption) *HTTPAPIClient {
	conf := DefaultHTTPConfig()
	for _, opt := range opts {
		opt(conf)
	}
	client := &http.Client{}
	return &HTTPAPIClient{
		addr:   addr,
		client: client,
		conf:   conf,
	}
}

func NewHTTPAPIClientWithClient(addr string, client httpClient, opts ...HTTPOption) *HTTPAPIClient {
	conf := DefaultHTTPConfig()
	for _, opt := range opts {
		opt(conf)
	}
	return &HTTPAPIClient{
		addr:   addr,
		client: client,
		conf:   conf,
	}
}

func (api *HTTPAPIClient) Get(path string) (*http.Response, error) {
	return api.fetch(http.MethodGet, path)
}

func (api *HTTPAPIClient) Post(path string) (*http.Response, error) {
	return api.fetch(http.MethodPost, path)
}

func (api *HTTPAPIClient) fetch(method, path string) (*http.Response, error) {
	url := &url.URL{
		Scheme: "http",
		Host:   api.addr,
		Path:   path,
	}
	req := &http.Request{
		Method: method,
		URL:    url,
		Header: api.conf.HTTPHeaders,
	}
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, WrapError(err, "failed to fetch")
	}
	return resp, nil
}

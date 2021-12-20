package gorqlite

import (
	"net/http"
)

type HTTPAPIClient struct {
	addr string
}

func NewHTTPAPIClient(addr string) *HTTPAPIClient {
	return &HTTPAPIClient{
		addr: addr,
	}
}

func (api *HTTPAPIClient) Get(path string) (*http.Response, error) {
	// TODO(AD) Check if path ends in '/'.
	return http.Get(api.addr + path)
}

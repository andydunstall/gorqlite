package gorqlite

import (
	"fmt"
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
	resp, err := http.Get(api.addr + path)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %s", err)
	}
	return resp, nil
}

//go:generate mockgen -source api.go -destination mocks/api/mock_api.go

package gorqlite

import (
	"context"
	"net/http"
	"net/url"
)

type APIClient interface {
	Get(path string, query url.Values) (*http.Response, error)
	GetWithContext(ctx context.Context, path string, query url.Values) (*http.Response, error)
	Post(path string, query url.Values, body []byte) (*http.Response, error)
	PostWithContext(ctx context.Context, path string, query url.Values, body []byte) (*http.Response, error)
}

func isStatusOK(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

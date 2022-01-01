//go:generate mockgen -source api.go -destination mocks/api/mock_api.go

package gorqlite

import (
	"context"
	"net/http"
)

type APIClient interface {
	Get(path string, opts ...Option) (*http.Response, error)
	GetWithContext(ctx context.Context, path string, opts ...Option) (*http.Response, error)
	Post(path string, body []byte, opts ...Option) (*http.Response, error)
	PostWithContext(ctx context.Context, path string, body []byte, opts ...Option) (*http.Response, error)
}

func isStatusOK(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

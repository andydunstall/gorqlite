//go:generate mockgen -source api.go -destination mocks/mock_api.go

package gorqlite

import (
	"context"
	"net/http"
)

type apiClient interface {
	Get(path string) (*http.Response, error)
	GetWithContext(ctx context.Context, path string) (*http.Response, error)
	Post(path string, body []byte) (*http.Response, error)
	PostWithContext(ctx context.Context, path string, body []byte) (*http.Response, error)
}

func isStatusOK(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

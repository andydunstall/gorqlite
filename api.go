//go:generate mockgen -source api.go -destination mocks/mock_api.go

package gorqlite

import (
	"net/http"
)

type APIClient interface {
	Get(path string) (*http.Response, error)
}

func IsStatusOK(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

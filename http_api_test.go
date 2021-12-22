package gorqlite

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/dunstall/gorqlite/mocks"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHTTPAPIClient_DefaultGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite",
			Path:   "/status",
		},
		Header: make(http.Header),
	}
	expectedResp := &http.Response{
		StatusCode: http.StatusOK,
	}
	httpClient := mock_gorqlite.NewMockhttpClient(ctrl)
	httpClient.EXPECT().Do(expectedReq).Return(expectedResp, nil)

	api := NewHTTPAPIClientWithClient("rqlite", httpClient)
	resp, err := api.Get("/status")
	require.Nil(t, err)
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_DefaultPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite",
			Path:   "/status",
		},
		Header: make(http.Header),
	}
	expectedResp := &http.Response{
		StatusCode: http.StatusOK,
	}
	httpClient := mock_gorqlite.NewMockhttpClient(ctrl)
	httpClient.EXPECT().Do(expectedReq).Return(expectedResp, nil)

	api := NewHTTPAPIClientWithClient("rqlite", httpClient)
	resp, err := api.Post("/status")
	require.Nil(t, err)
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_GetWithConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	headers := make(http.Header)
	headers["abc"] = []string{"xyz"}

	expectedReq := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite",
			Path:   "/status",
		},
		Header: headers,
	}
	expectedResp := &http.Response{
		StatusCode: http.StatusOK,
	}
	httpClient := mock_gorqlite.NewMockhttpClient(ctrl)
	httpClient.EXPECT().Do(expectedReq).Return(expectedResp, nil)

	api := NewHTTPAPIClientWithClient(
		"rqlite",
		httpClient,
		WithHTTPHeaders(headers),
	)
	resp, err := api.Get("/status")
	require.Nil(t, err)
	require.Equal(t, expectedResp, resp)
}

package gorqlite

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "rqlite",
	}
	expectedResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	httpClient := mock_gorqlite.NewMockhttpClient(ctrl)
	httpClient.EXPECT().Do(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := NewHTTPAPIClientWithClient("rqlite", httpClient)
	resp, err := api.Get("/status")
	require.Nil(t, err)
	defer resp.Body.Close()
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
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "rqlite",
	}
	expectedResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	httpClient := mock_gorqlite.NewMockhttpClient(ctrl)
	httpClient.EXPECT().Do(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := NewHTTPAPIClientWithClient("rqlite", httpClient)
	resp, err := api.Post("/status", nil)
	require.Nil(t, err)
	defer resp.Body.Close()
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
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     headers,
		Host:       "rqlite",
	}
	expectedResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	httpClient := mock_gorqlite.NewMockhttpClient(ctrl)
	httpClient.EXPECT().Do(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := NewHTTPAPIClientWithClient(
		"rqlite",
		httpClient,
		WithHTTPHeaders(headers),
	)
	resp, err := api.Get("/status")
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

type httpReqEqMatcher struct {
	x interface{}
}

func newHTTPReqEqMatcher(r *http.Request) gomock.Matcher {
	return &httpReqEqMatcher{
		x: r,
	}
}

func (e httpReqEqMatcher) Matches(x interface{}) bool {
	lhs, ok := e.x.(*http.Request)
	if !ok {
		return false
	}
	rhs, ok := x.(*http.Request)
	if !ok {
		return false
	}

	// Removes unexported fields to compare.
	strippedLHS := &http.Request{
		Method: lhs.Method,
		URL:    lhs.URL,
		Proto:  lhs.Proto,
		Header: lhs.Header,
		Host:   lhs.Host,
		Body:   lhs.Body,
	}
	strippedRHS := &http.Request{
		Method: rhs.Method,
		URL:    rhs.URL,
		Proto:  rhs.Proto,
		Header: rhs.Header,
		Host:   rhs.Host,
		Body:   rhs.Body,
	}
	return reflect.DeepEqual(strippedLHS, strippedRHS)
}

func (e httpReqEqMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.x)
}

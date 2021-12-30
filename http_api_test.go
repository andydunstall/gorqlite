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

var (
	testAddrs = []string{"rqlite"}
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
	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := NewHTTPAPIClient(testAddrs, WithTransport(transport))
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
	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := NewHTTPAPIClient(testAddrs, WithTransport(transport))
	resp, err := api.Post("/status", nil)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_GetWithHeaders(t *testing.T) {
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
	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := NewHTTPAPIClient(
		testAddrs,
		WithHTTPHeaders(headers),
		WithTransport(transport),
	)
	resp, err := api.Get("/status")
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_WithActiveHostRoundRobin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite-0", "rqlite-1", "rqlite-2"}

	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	api := NewHTTPAPIClient(addrs, WithTransport(transport))

	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			expectedReq := &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   addrs[j],
					Path:   "/status",
				},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     make(http.Header),
				Host:       addrs[j],
			}
			expectedResp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}
			transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

			resp, err := api.Get("/status")
			require.Nil(t, err)
			defer resp.Body.Close()
			require.Equal(t, expectedResp, resp)
		}
	}
}

func TestHTTPAPIClient_WithoutActiveHostRoundRobin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite", "rqlite-0", "rqlite-0"}

	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	api := NewHTTPAPIClient(addrs, WithTransport(transport), WithActiveHostRoundRobin(false))

	for i := 0; i < 4; i++ {
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
		transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

		resp, err := api.Get("/status")
		require.Nil(t, err)
		defer resp.Body.Close()
		require.Equal(t, expectedResp, resp)
	}
}

func TestHTTPAPIClient_WithRedirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	// Update the redirect path for each request to check its being updated.
	for i := 0; i < 4; i++ {
		expectedReq := &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: "http",
				Host:   "rqlite",
				Path:   fmt.Sprintf("/status-%d", i),
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Host:       "rqlite",
		}
		resp := &http.Response{
			StatusCode: http.StatusMovedPermanently,
			Header: http.Header{
				"Location": []string{fmt.Sprintf("http://rqlite/status-%d", i+1)},
			},
			Body: ioutil.NopCloser(strings.NewReader("")),
		}
		transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(resp, nil)
	}

	api := NewHTTPAPIClient(testAddrs, WithRedirectAttempts(3), WithTransport(transport))
	resp, err := api.Get("/status-0")
	require.Error(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
}

func TestHTTPAPIClient_NoRedirects(t *testing.T) {
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
		StatusCode: http.StatusMovedPermanently,
		Header: http.Header{
			"Location": []string{"http://0.0.0.0:0"},
		},
		Body: ioutil.NopCloser(strings.NewReader("")),
	}
	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := NewHTTPAPIClient(testAddrs, WithTransport(transport), WithRedirectAttempts(0))
	resp, err := api.Get("/status")
	require.Error(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
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

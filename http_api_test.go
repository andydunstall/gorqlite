package gorqlite_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dunstall/gorqlite"
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

	api := gorqlite.NewHTTPAPIClient(testAddrs, gorqlite.WithTransport(transport))
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

	api := gorqlite.NewHTTPAPIClient(testAddrs, gorqlite.WithTransport(transport))
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

	api := gorqlite.NewHTTPAPIClient(
		testAddrs,
		gorqlite.WithHTTPHeaders(headers),
		gorqlite.WithTransport(transport),
	)
	resp, err := api.Get("/status")
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_RetryFailedRequestsSucceeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite-badstatus", "rqlite-network", "rqlite-ok"}

	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	clock := mock_gorqlite.NewMockClock(ctrl)
	api := gorqlite.NewHTTPAPIClient(
		addrs,
		gorqlite.WithTransport(transport),
		gorqlite.WithClock(clock),
		// Disable to check still tries all nodes.
		gorqlite.WithActiveHostRoundRobin(false),
	)

	clock.EXPECT().Sleep(100 * time.Millisecond)
	clock.EXPECT().Sleep(200 * time.Millisecond)

	// First return a bad status.
	expectedReq1 := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite-badstatus",
			Path:   "/status",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "rqlite-badstatus",
	}
	expectedResp1 := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq1)).Return(expectedResp1, nil)

	// Next return a network error.
	expectedReq2 := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite-network",
			Path:   "/status",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "rqlite-network",
	}
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq2)).Return(nil, fmt.Errorf("network error"))

	// Return OK from the final host.
	expectedReq3 := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite-ok",
			Path:   "/status",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "rqlite-ok",
	}
	expectedResp3 := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq3)).Return(expectedResp3, nil)

	resp, err := api.Get("/status")
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp3, resp)
}

func TestHTTPAPIClient_RetryFailedRequestsFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite-badstatus", "rqlite-network"}

	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	clock := mock_gorqlite.NewMockClock(ctrl)
	api := gorqlite.NewHTTPAPIClient(addrs, gorqlite.WithTransport(transport), gorqlite.WithClock(clock))

	for i := 0; i < 6; i++ {
		d := (100 << i) * time.Millisecond
		clock.EXPECT().Sleep(d)
	}

	// First return a bad status.
	expectedReq1 := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite-badstatus",
			Path:   "/status",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "rqlite-badstatus",
	}
	expectedResp1 := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}

	// Next return a network error.
	expectedReq2 := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "rqlite-network",
			Path:   "/status",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "rqlite-network",
	}

	// Expect 6 retries.
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq1)).Return(expectedResp1, nil)
	for i := 0; i < 3; i++ {
		transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq2)).Return(nil, fmt.Errorf("network error"))
		transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq1)).Return(expectedResp1, nil)
	}

	resp, err := api.Get("/status")
	require.Error(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
}

func TestHTTPAPIClient_FailureNotRetryable(t *testing.T) {
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
		StatusCode: http.StatusForbidden,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := gorqlite.NewHTTPAPIClient(testAddrs, gorqlite.WithTransport(transport))
	resp, err := api.Get("/status")
	require.Error(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
}

func TestHTTPAPIClient_WithActiveHostRoundRobin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite-0", "rqlite-1", "rqlite-2"}

	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	api := gorqlite.NewHTTPAPIClient(addrs, gorqlite.WithTransport(transport))

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

	addrs := []string{"rqlite", "rqlite-0", "rqlite-1"}

	transport := mock_gorqlite.NewMockRoundTripper(ctrl)
	api := gorqlite.NewHTTPAPIClient(addrs, gorqlite.WithTransport(transport), gorqlite.WithActiveHostRoundRobin(false))

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

func TestHTTPAPIClient_Transaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "http",
			Host:     "rqlite",
			Path:     "/status",
			RawQuery: "transaction=",
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

	api := gorqlite.NewHTTPAPIClient(testAddrs, gorqlite.WithTransport(transport))
	resp, err := api.Get("/status", gorqlite.WithTransaction(true))
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_Consistency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "http",
			Host:     "rqlite",
			Path:     "/status",
			RawQuery: "level=strong",
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

	api := gorqlite.NewHTTPAPIClient(testAddrs, gorqlite.WithTransport(transport))
	resp, err := api.Get("/status", gorqlite.WithConsistency("strong"))
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_ConfigOverride(t *testing.T) {
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
	// Ignore the first request and verify the second has no query string.
	transport.EXPECT().RoundTrip(gomock.Any()).Return(expectedResp, nil)
	transport.EXPECT().RoundTrip(newHTTPReqEqMatcher(expectedReq)).Return(expectedResp, nil)

	api := gorqlite.NewHTTPAPIClient(testAddrs, gorqlite.WithTransport(transport))
	resp, err := api.Get("/status", gorqlite.WithConsistency("strong"))
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
	resp, err = api.Get("/status")
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

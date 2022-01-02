package gorqlite

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dunstall/gorqlite/mocks/http_api"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	testAddrs = []string{"rqlite"}
)

func TestHTTPAPIClient_DefaultGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq, err := http.NewRequest(http.MethodGet, "http://rqlite/status", nil)
	require.Nil(t, err)
	expectedResp := httpResponse(http.StatusOK, strings.NewReader(""))
	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq),
	).Return(expectedResp, nil)

	api := newHTTPAPIClient(testAddrs, transport, &systemClock{}, true)
	resp, err := api.Get("/status", url.Values{})
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_DefaultPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq, err := http.NewRequest(http.MethodPost, "http://rqlite/status", nil)
	require.Nil(t, err)
	expectedResp := httpResponse(http.StatusOK, strings.NewReader(""))
	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq),
	).Return(expectedResp, nil)

	api := newHTTPAPIClient(testAddrs, transport, &systemClock{}, true)
	resp, err := api.Post("/status", url.Values{}, nil)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_FetchWithQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq, err := http.NewRequest(http.MethodGet, "http://rqlite/status?a=b", nil)
	require.Nil(t, err)
	expectedResp := httpResponse(http.StatusOK, strings.NewReader(""))
	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq),
	).Return(expectedResp, nil)

	query := url.Values{}
	query.Add("a", "b")

	api := newHTTPAPIClient(testAddrs, transport, &systemClock{}, true)
	resp, err := api.Get("/status", query)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp, resp)
}

func TestHTTPAPIClient_RetryFailedRequestsSucceeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite-badstatus", "rqlite-network", "rqlite-ok"}

	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	clock := mock_gorqlite.NewMockclock(ctrl)
	// Disable round robin to check still tries all nodes.
	api := newHTTPAPIClient(addrs, transport, clock, false)

	clock.EXPECT().Sleep(100 * time.Millisecond)
	clock.EXPECT().Sleep(200 * time.Millisecond)

	// First return a bad status.
	expectedReq1, err := http.NewRequest(
		http.MethodGet, "http://rqlite-badstatus/status", nil,
	)
	require.Nil(t, err)
	expectedResp1 := httpResponse(
		http.StatusInternalServerError, strings.NewReader(""),
	)
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq1),
	).Return(expectedResp1, nil)

	// Next return a network error.
	expectedReq2, err := http.NewRequest(
		http.MethodGet, "http://rqlite-network/status", nil,
	)
	require.Nil(t, err)
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq2),
	).Return(nil, fmt.Errorf("network error"))

	// Return OK from the final host.
	expectedReq3, err := http.NewRequest(
		http.MethodGet, "http://rqlite-ok/status", nil,
	)
	require.Nil(t, err)
	expectedResp3 := httpResponse(http.StatusOK, strings.NewReader(""))
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq3),
	).Return(expectedResp3, nil)

	resp, err := api.Get("/status", url.Values{})
	require.Nil(t, err)
	defer resp.Body.Close()
	require.Equal(t, expectedResp3, resp)
}

func TestHTTPAPIClient_RetryFailedRequestsFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite-badstatus", "rqlite-network"}

	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	clock := mock_gorqlite.NewMockclock(ctrl)
	api := newHTTPAPIClient(addrs, transport, clock, true)

	for i := 0; i < 6; i++ {
		d := (100 << i) * time.Millisecond
		clock.EXPECT().Sleep(d)
	}

	// First return a bad status.
	expectedReq1, err := http.NewRequest(
		http.MethodGet, "http://rqlite-badstatus/status", nil,
	)
	require.Nil(t, err)
	expectedResp1 := httpResponse(
		http.StatusInternalServerError, strings.NewReader(""),
	)

	// Next return a network error.
	expectedReq2, err := http.NewRequest(
		http.MethodGet, "http://rqlite-network/status", nil,
	)
	require.Nil(t, err)

	// Expect 6 retries.
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq1),
	).Return(expectedResp1, nil)
	for i := 0; i < 3; i++ {
		transport.EXPECT().RoundTrip(
			newHTTPReqEqMatcher(expectedReq2),
		).Return(nil, fmt.Errorf("network error"))
		transport.EXPECT().RoundTrip(
			newHTTPReqEqMatcher(expectedReq1),
		).Return(expectedResp1, nil)
	}

	_, err = api.Get("/status", url.Values{})
	require.Error(t, err)
}

func TestHTTPAPIClient_FailureNotRetryable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedReq, err := http.NewRequest(http.MethodGet, "http://rqlite/status", nil)
	require.Nil(t, err)
	expectedResp := httpResponse(http.StatusForbidden, strings.NewReader(""))
	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	transport.EXPECT().RoundTrip(
		newHTTPReqEqMatcher(expectedReq),
	).Return(expectedResp, nil)

	api := newHTTPAPIClient(testAddrs, transport, &systemClock{}, true)
	_, err = api.Get("/status", url.Values{})
	require.Error(t, err)
}

func TestHTTPAPIClient_WithActiveHostRoundRobin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{"rqlite-0", "rqlite-1", "rqlite-2"}

	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	api := newHTTPAPIClient(addrs, transport, &systemClock{}, true)

	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			expectedReq, err := http.NewRequest(
				http.MethodGet, fmt.Sprintf("http://%s/status", addrs[j]), nil,
			)
			require.Nil(t, err)
			expectedResp := httpResponse(http.StatusOK, strings.NewReader(""))
			transport.EXPECT().RoundTrip(
				newHTTPReqEqMatcher(expectedReq),
			).Return(expectedResp, nil)

			resp, err := api.Get("/status", url.Values{})
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

	transport := mock_gorqlite.NewMockroundTripper(ctrl)
	api := newHTTPAPIClient(addrs, transport, &systemClock{}, false)

	for i := 0; i < 4; i++ {
		expectedReq, err := http.NewRequest(http.MethodGet, "http://rqlite/status", nil)
		require.Nil(t, err)
		expectedResp := httpResponse(http.StatusOK, strings.NewReader(""))
		transport.EXPECT().RoundTrip(
			newHTTPReqEqMatcher(expectedReq),
		).Return(expectedResp, nil)

		resp, err := api.Get("/status", url.Values{})
		require.Nil(t, err)
		defer resp.Body.Close()
		require.Equal(t, expectedResp, resp)
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

func httpResponse(statusCode int, body io.Reader) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(body),
	}
}

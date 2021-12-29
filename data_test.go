package gorqlite

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/dunstall/gorqlite/mocks"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestDataAPIClient_QueryOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/query", []byte(`["abc","123"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	err := dataClient.Query([]string{"abc", "123"})
	require.Nil(t, err)
}

func TestDataAPIClient_QueryBadStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resp := &http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/query", []byte(`["abc","123"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	err := dataClient.Query([]string{"abc", "123"})
	require.Error(t, err)
}

func TestDataAPIClient_QueryNetworkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/query", []byte(`["abc","123"]`)).Return(nil, fmt.Errorf("network err"))

	dataClient := NewDataAPIClientWithClient(apiClient)
	err := dataClient.Query([]string{"abc", "123"})
	require.Error(t, err)
}

func TestDataAPIClient_ExecuteOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/execute", []byte(`["abc","123"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	err := dataClient.Execute([]string{"abc", "123"})
	require.Nil(t, err)
}

func TestDataAPIClient_ExecuteBadStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resp := &http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/execute", []byte(`["abc","123"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	err := dataClient.Execute([]string{"abc", "123"})
	require.Error(t, err)
}

func TestDataAPIClient_ExecuteNetworkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/execute", []byte(`["abc","123"]`)).Return(nil, fmt.Errorf("network err"))

	dataClient := NewDataAPIClientWithClient(apiClient)
	err := dataClient.Execute([]string{"abc", "123"})
	require.Error(t, err)
}

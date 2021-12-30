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

// Note test data from https://github.com/rqlite/gorqlite/blob/master/query.go.

func TestDataAPIClient_QueryOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{
    "results": [
        {
            "columns": [
                "id",
                "name"
            ],
            "types": [
                "integer",
                "text"
            ],
            "values": [
                [
                    1,
                    "foo"
                ],
                [
                    2,
                    "bar"
                ]
            ],
            "time": 10
        }
    ],
    "time": 100
}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/query", []byte(`["SELECT * FROM mytable"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	result, err := dataClient.Query([]string{"SELECT * FROM mytable"})
	require.Nil(t, err)

	expectedResult := QueryResponse{
		Results: []QueryRows{
			{
				Columns: []string{"id", "name"},
				Types:   []string{"integer", "text"},
				Values: [][]interface{}{
					{
						float64(1), "foo",
					},
					{
						float64(2), "bar",
					},
				},
				Time: 10,
			},
		},
		Error: "",
		Time:  100,
	}
	require.Equal(t, expectedResult, result)
}

func TestDataAPIClient_QueryNullResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{
    "results": [
        {
            "columns": [
                "id",
                "name"
            ],
            "types": [
                "number",
                "text"
            ],
            "values": [
                [
                    null,
                    "Hulk"
                ]
            ],
            "time": 4
        },
        {
            "columns": [
                "id",
                "name"
            ],
            "types": [
                "number",
                "text"
            ],
            "time": 1
        }
    ],
    "time": 3
}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/query", []byte(`["SELECT * FROM mytable"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	result, err := dataClient.Query([]string{"SELECT * FROM mytable"})
	require.Nil(t, err)

	expectedResult := QueryResponse{
		Results: []QueryRows{
			{
				Columns: []string{"id", "name"},
				Types:   []string{"number", "text"},
				Values: [][]interface{}{
					{
						nil, "Hulk",
					},
				},
				Time: 4,
			},
			{
				Columns: []string{"id", "name"},
				Types:   []string{"number", "text"},
				Time:    1,
			},
		},
		Error: "",
		Time:  3,
	}
	require.Equal(t, expectedResult, result)
}

func TestDataAPIClient_QueryErrorResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{
    "results": [
        {
            "error": "near \"invalid\": syntax error"
        }
    ],
    "time": 2
}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/query", []byte(`["invalid"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	result, err := dataClient.Query([]string{"invalid"})
	require.Nil(t, err)

	expectedResult := QueryResponse{
		Results: []QueryRows{
			{
				Error: "near \"invalid\": syntax error",
			},
		},
		Time: 2,
	}
	require.Equal(t, expectedResult, result)
	require.Equal(t, "near \"invalid\": syntax error", result.GetFirstError())
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
	_, err := dataClient.Query([]string{"abc", "123"})
	require.Error(t, err)
}

func TestDataAPIClient_QueryNetworkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/query", []byte(`["abc","123"]`)).Return(nil, fmt.Errorf("network err"))

	dataClient := NewDataAPIClientWithClient(apiClient)
	_, err := dataClient.Query([]string{"abc", "123"})
	require.Error(t, err)
}

func TestDataAPIClient_ExecuteOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{
    "results": [
        {
            "last_insert_id": 1,
            "rows_affected": 1,
            "time": 10
        },
        {
            "last_insert_id": 2,
            "rows_affected": 1,
            "time": 20
        }
    ],
    "time": 100
}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/execute", []byte(`["abc","123"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	result, err := dataClient.Execute([]string{"abc", "123"})
	require.Nil(t, err)

	expectedResult := ExecuteResponse{
		Results: []ExecuteResult{
			{
				LastInsertId: 1,
				RowsAffected: 1,
				Time:         10,
			},
			{
				LastInsertId: 2,
				RowsAffected: 1,
				Time:         20,
			},
		},
		Error: "",
		Time:  100,
	}
	require.Equal(t, expectedResult, result)
}

func TestDataAPIClient_ExecuteErrorResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{
    "results": [
        {
            "last_insert_id": 1,
            "rows_affected": 1,
            "time": 10
        },
        {
            "error": "invalid request"
        }
    ],
    "time": 100
}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/execute", []byte(`["abc","123"]`)).Return(resp, nil)

	dataClient := NewDataAPIClientWithClient(apiClient)
	result, err := dataClient.Execute([]string{"abc", "123"})
	require.Nil(t, err)

	expectedResult := ExecuteResponse{
		Results: []ExecuteResult{
			{
				LastInsertId: 1,
				RowsAffected: 1,
				Time:         10,
			},
			{
				Error: "invalid request",
			},
		},
		Error: "",
		Time:  100,
	}
	require.Equal(t, expectedResult, result)
	require.Equal(t, "invalid request", result.GetFirstError())
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
	_, err := dataClient.Execute([]string{"abc", "123"})
	require.Error(t, err)
}

func TestDataAPIClient_ExecuteNetworkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClient := mock_gorqlite.NewMockAPIClient(ctrl)
	apiClient.EXPECT().PostWithContext(gomock.Any(), "/db/execute", []byte(`["abc","123"]`)).Return(nil, fmt.Errorf("network err"))

	dataClient := NewDataAPIClientWithClient(apiClient)
	_, err := dataClient.Execute([]string{"abc", "123"})
	require.Error(t, err)
}

package gorqlite

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueryRow_ScanToTime(t *testing.T) {
	row := QueryRow{
		Columns: []string{
			"int64-to-time",
			"float64-to-time",
			"string-to-time",
		},
		Values: []interface{}{
			int64(1641067851),
			float64(1641067851.45),
			"2022-01-01T20:10:51+00:00",
		},
	}
	var col1 time.Time
	var col2 time.Time
	var col3 time.Time
	require.Nil(t, row.Scan(&col1, &col2, &col3))
	expectedTime := time.Time(time.Date(2022, time.January, 1, 20, 10, 51, 0, time.Local))
	require.Equal(t, expectedTime, col1)
	require.Equal(t, expectedTime, col2)
	require.Equal(t, expectedTime, col3)
}

func TestQueryRow_ScanToFloat64(t *testing.T) {
	row := QueryRow{
		Columns: []string{
			"int64-to-float64",
			"float64-to-float64",
			"string-to-float64",
		},
		Values: []interface{}{
			int64(123),
			float64(123.45),
			"123.45",
		},
	}
	var col1 float64
	var col2 float64
	var col3 float64
	require.Nil(t, row.Scan(&col1, &col2, &col3))
	require.Equal(t, float64(123), col1)
	require.Equal(t, float64(123.45), col2)
	require.Equal(t, float64(123.45), col3)
}

func TestQueryRow_ScanToInt64(t *testing.T) {
	row := QueryRow{
		Columns: []string{
			"int64-to-int64",
			"float64-to-int64",
			"string-to-int64",
		},
		Values: []interface{}{
			int64(123),
			float64(123.45),
			"123",
		},
	}
	var col1 int64
	var col2 int64
	var col3 int64
	require.Nil(t, row.Scan(&col1, &col2, &col3))
	require.Equal(t, int64(123), col1)
	require.Equal(t, int64(123), col2)
	require.Equal(t, int64(123), col3)
}

func TestQueryRow_ScanToInt(t *testing.T) {
	row := QueryRow{
		Columns: []string{
			"int64-to-int",
			"float64-to-int",
			"string-to-int",
		},
		Values: []interface{}{
			int64(123),
			float64(123.45),
			"123",
		},
	}
	var col1 int
	var col2 int
	var col3 int
	require.Nil(t, row.Scan(&col1, &col2, &col3))
	require.Equal(t, 123, col1)
	require.Equal(t, 123, col2)
	require.Equal(t, 123, col3)
}

func TestQueryRow_ScanToString(t *testing.T) {
	row := QueryRow{
		Columns: []string{
			"string-to-string",
		},
		Values: []interface{}{
			"mystring",
		},
	}
	var col1 string
	require.Nil(t, row.Scan(&col1))
	require.Equal(t, "mystring", col1)
}

func TestQueryRow_ScanSkipNillValues(t *testing.T) {
	row := QueryRow{
		Columns: []string{"1", "2"},
		Values:  []interface{}{int64(1), nil},
	}
	var a int
	var b int
	require.Nil(t, row.Scan(&a, &b))
	require.Equal(t, 1, a)
	require.Equal(t, 0, b)
}

func TestQueryRow_ScanInsufficientVars(t *testing.T) {
	row := QueryRow{
		Columns: []string{"1", "2"},
		Values:  []interface{}{1, 2},
	}
	var a int
	require.Error(t, row.Scan(&a))
}

func TestQueryRow_ScanInsufficientValues(t *testing.T) {
	row := QueryRow{
		Columns: []string{"1", "2"},
		Values:  []interface{}{1},
	}
	var a int
	var b int
	require.Error(t, row.Scan(&a, &b))
}

package gorqlite

import (
	"strconv"
	"time"
)

type QueryRow struct {
	Columns []string
	Values  []interface{}
}

// Based on https://github.com/rqlite/gorqlite/blob/5cf06496fee7a89243002b194bd095eba64346b3/query.go#L385.
func (r *QueryRow) Scan(vars ...interface{}) error {
	if len(r.Columns) != len(r.Values) {
		return newError(
			"invalid row: incorrect number of types, was %d, needed %d",
			len(r.Values), len(r.Columns),
		)
	}
	if len(vars) != len(r.Columns) {
		return newError(
			"invalid number of vars: was %d, needed %d",
			len(vars), len(r.Columns),
		)
	}

	for i, dest := range vars {
		src := r.Values[i]
		// Skip nil items.
		if src == nil {
			continue
		}

		switch dest.(type) {
		case *string:
			switch src := src.(type) {
			case string:
				*dest.(*string) = src
			default:
				return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
			}
		case *int:
			switch src := src.(type) {
			case int64:
				*dest.(*int) = int(src)
			case float64:
				*dest.(*int) = int(src)
			case string:
				n, err := strconv.Atoi(src)
				if err != nil {
					return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
				}
				*dest.(*int) = n
			default:
				return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
			}
		case *int64:
			switch src := src.(type) {
			case int64:
				*dest.(*int64) = src
			case float64:
				*dest.(*int64) = int64(src)
			case string:
				n, err := strconv.ParseInt(src, 10, 64)
				if err != nil {
					return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
				}
				*dest.(*int64) = n
			default:
				return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
			}
		case *float64:
			switch src := src.(type) {
			case float64:
				*dest.(*float64) = src
			case int64:
				*dest.(*float64) = float64(src)
			case string:
				n, err := strconv.ParseFloat(src, 64)
				if err != nil {
					return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
				}
				*dest.(*float64) = n
			default:
				return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
			}
		case *time.Time:
			t, err := toTime(src)
			if err != nil {
				return newError("invalid conversion from %T to %T (value %v)", src, dest, dest)
			}
			*dest.(*time.Time) = t
		default:
			return newError("unsupported destination type: %T", dest)
		}
	}

	return nil
}

type QueryResult struct {
	Columns []string        `json:"columns,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
	Error   string          `json:"error,omitempty"`
	row     int
}

func (r *QueryResult) Rows() int {
	return len(r.Values)
}

func (r *QueryResult) Next() (*QueryRow, bool) {
	if r.row >= r.Rows() {
		return nil, false
	}
	row := &QueryRow{
		Columns: r.Columns,
		Values:  r.Values[r.row],
	}
	r.row++
	return row, true
}

type QueryResults []QueryResult

func (r QueryResults) GetFirstError() string {
	for _, row := range r {
		if row.Error != "" {
			return row.Error
		}
	}
	return ""
}

func (r QueryResults) HasError() bool {
	return r.GetFirstError() != ""
}

type ExecuteResult struct {
	LastInsertId int64  `json:"last_insert_id,omitempty"`
	RowsAffected int64  `json:"rows_affected,omitempty"`
	Error        string `json:"error,omitempty"`
}

type ExecuteResults []ExecuteResult

func (r ExecuteResults) GetFirstError() string {
	for _, result := range r {
		if result.Error != "" {
			return result.Error
		}
	}
	return ""
}

func (r ExecuteResults) HasError() bool {
	return r.GetFirstError() != ""
}

type Nodes map[string]struct {
	APIAddr   string  `json:"api_addr,omitempty"`
	Addr      string  `json:"addr,omitempty"`
	Reachable bool    `json:"reachable"`
	Leader    bool    `json:"leader"`
	Time      float64 `json:"time,omitempty"`
	Error     string  `json:"error,omitempty"`
}

func toTime(src interface{}) (time.Time, error) {
	switch src := src.(type) {
	case int64:
		return time.Unix(src, 0), nil
	case float64:
		return time.Unix(int64(src), 0), nil
	case string:
		return time.Parse(time.RFC3339, src)
	default:
		return time.Time{}, newError("invalid time convertion from %T (value: %v)", src, src)
	}
}

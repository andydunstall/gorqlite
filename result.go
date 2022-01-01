package gorqlite

type QueryResult struct {
	Columns []string        `json:"columns,omitempty"`
	Types   []string        `json:"types,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
	Error   string          `json:"error,omitempty"`
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

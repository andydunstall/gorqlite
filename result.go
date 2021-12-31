package gorqlite

type QueryRows struct {
	Columns []string        `json:"columns,omitempty"`
	Types   []string        `json:"types,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type QueryResponse struct {
	Results []QueryRows `json:"results,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (r *QueryResponse) GetFirstError() string {
	if r.Error != "" {
		return r.Error
	}
	for _, row := range r.Results {
		if row.Error != "" {
			return row.Error
		}
	}
	return ""
}

func (r *QueryResponse) HasError() bool {
	return r.GetFirstError() != ""
}

type ExecuteResult struct {
	LastInsertId int64  `json:"last_insert_id,omitempty"`
	RowsAffected int64  `json:"rows_affected,omitempty"`
	Error        string `json:"error,omitempty"`
}

type ExecuteResponse struct {
	Results []ExecuteResult `json:"results,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func (r *ExecuteResponse) GetFirstError() string {
	if r.Error != "" {
		return r.Error
	}
	for _, result := range r.Results {
		if result.Error != "" {
			return result.Error
		}
	}
	return ""
}

func (r *ExecuteResponse) HasError() bool {
	return r.GetFirstError() != ""
}

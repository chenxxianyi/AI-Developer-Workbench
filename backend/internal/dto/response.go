package dto

// PaginatedResponse is the generic paginated response matching frontend PaginatedData<T>.
type PaginatedResponse[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// ListReportsQuery holds query parameters for listing reports.
type ListReportsQuery struct {
	ToolType string `form:"tool_type"`
	Status   string `form:"status"`
	Sort     string `form:"sort"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// SetDefaults applies default values for list query.
func (q *ListReportsQuery) SetDefaults() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
}

// ValidSortValues returns allowed sort values.
func ValidSortValues() map[string]string {
	return map[string]string{
		"newest":     "created_at DESC",
		"oldest":     "created_at ASC",
		"score_desc": "total_score DESC",
		"score_asc":  "total_score ASC",
	}
}

// ValidStatusValues returns the set of allowed report status filters.
func ValidStatusValues() map[string]bool {
	return map[string]bool{
		"processing": true,
		"succeeded":  true,
		"fallback":   true,
		"failed":     true,
	}
}

package api

type PaginatedResponse struct {
	Data       []any `json:"data"`
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalCount int64 `json:"total_count"`
}

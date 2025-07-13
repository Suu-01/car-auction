package api

// PaginatedResponse는 페이징된 목록 응답에 공통으로 사용됩니다.
type PaginatedResponse struct {
	Data       []any `json:"data"`
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalCount int64 `json:"total_count"`
}

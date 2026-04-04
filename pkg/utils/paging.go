package utils

import "math"

const (
	DEFAULT_PAGESIZE int64 = 10
)

type Pagination struct {
	CurrentPage int64 `json:"current_page"`
	Total       int64 `json:"total"`
	TotalPage   int64 `json:"total_page"`
	Limit       int64 `json:"limit"`
	Skip        int64 `json:"skip"`
}

func NewPaging(page int64, pageSize int64, total int64) *Pagination {
	var pageInfo Pagination
	limit := DEFAULT_PAGESIZE
	if pageSize > 0 && pageSize <= limit {
		pageInfo.Limit = pageSize
	} else {
		pageInfo.Limit = DEFAULT_PAGESIZE
	}
	totalPage := int64(math.Ceil(float64(total) / float64(pageInfo.Limit)))
	pageInfo.Total = total
	pageInfo.TotalPage = totalPage
	if page <= 0 || totalPage == 1 {
		page = 1
	}
	pageInfo.CurrentPage = page
	pageInfo.Skip = (page - 1) * pageSize

	return &pageInfo
}

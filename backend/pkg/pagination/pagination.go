package pagination

type Pagination struct {
	CurrentPage     int64   `json:"current_page"`
	CurrentElements int64   `json:"current_elements"`
	TotalPages      int64   `json:"total_pages"`
	TotalElements   int64   `json:"total_elements"`
	SortBy          string  `json:"sort_by"`
	CursorStart     *string `json:"cursor_start,omitempty"`
	CursorEnd       *string `json:"cursor_end,omitempty"`
}

type PaginationRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func (p *PaginationRequest) GetPageWithDefault() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PaginationRequest) GetPageSizeWithDefault() int {
	if p.PageSize <= 0 {
		return 20
	}
	return p.PageSize
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPageWithDefault() - 1) * p.GetPageSizeWithDefault()
}

func (p *PaginationRequest) GetLimit() int {
	return p.GetPageSizeWithDefault()
}

func ResetPagination() *Pagination {

	return &Pagination{
		CurrentPage:     1,
		CurrentElements: 0,
		TotalPages:      0,
		TotalElements:   0,
	}
}

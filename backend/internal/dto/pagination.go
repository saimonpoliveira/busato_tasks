package dto

type PaginationQuery struct {
	Page      int    `form:"page" json:"page"`
	PageSize  int    `form:"page_size" json:"page_size"`
	SortBy    string `form:"sort_by" json:"sort_by"`
	SortOrder string `form:"sort_order" json:"sort_order"`
	Search    string `form:"search" json:"search"`
}

func (p *PaginationQuery) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.SortOrder != "asc" && p.SortOrder != "desc" {
		p.SortOrder = "desc"
	}
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginatedResponse[T any](data []T, total int64, page, pageSize int) PaginatedResponse[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	if data == nil {
		data = []T{}
	}
	return PaginatedResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

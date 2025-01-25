package helper

import "math"

type Paginator struct {
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	TotalRows   int64    `json:"total_rows"`
	TotalPages  int      `json:"total_pages"`
	SortColumns []string `json:"sort_columns,omitempty"`
	SortOrders  []string `json:"sort_orders,omitempty"`
	IsFirstPage bool     `json:"is_first_page"`
	IsLastPage  bool     `json:"is_last_page"`
}

func NewPaginator(page, limit int) *Paginator {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return &Paginator{
		Page:  page,
		Limit: limit,
	}
}

func (p *Paginator) Offset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

func (p *Paginator) CalculateTotalPages() {
	if p.Limit > 0 {
		p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.Limit)))
	} else {
		p.TotalPages = 0
	}
	p.IsFirstPage = (p.Page == 1)
	p.IsLastPage = (p.Page >= p.TotalPages)
}

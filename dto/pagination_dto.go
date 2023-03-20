package dto

type Pagination struct {
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
	Page      int   `json:"page,omitempty"`
}

type PaginationRequest struct {
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
	OrderBy string `json:"orderBy"`
	Order   string `json:"order"`
	Search  string `json:"search"`
}

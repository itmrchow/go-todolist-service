package entity

// Pagination represents pagination and sorting information
type Pagination struct {
	Page      int    `json:"page" binding:"min=1"`
	PageSize  int    `json:"page_size" binding:"min=1,max=100"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order" binding:"oneof=asc desc"`
}

// PaginationResponse represents pagination information in the response
type PaginationResponse struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalCount  int  `json:"total_count"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

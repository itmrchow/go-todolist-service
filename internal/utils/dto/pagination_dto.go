package dto

// PaginationReq represents pagination and sorting information
type PaginationReq struct {
	Page      int    `json:"page" binding:"min=1"`
	PageSize  int    `json:"page_size" binding:"min=1,max=100"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order" binding:"oneof=asc desc"`
}

// PaginationResp represents pagination information in the response
type PaginationResp struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

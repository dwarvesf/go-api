package viewmodel

// Pagination is the response for pagination
type Pagination struct {
	Page         int `json:"page" validate:"required"`
	PageSize     int `json:"page_size" validate:"required"`
	TotalPages   int `json:"total_pages" validate:"required"`
	TotalRecords int `json:"total_records" validate:"required"`
} // @name Pagination

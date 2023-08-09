package model

// Pagination represent the pagination
type Pagination struct {
	Page         int
	PageSize     int
	TotalRecords int
	TotalPages   int
	Offset       int
}

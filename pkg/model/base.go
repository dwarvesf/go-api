package model

// ListResult represent the list result
type ListResult[T any] struct {
	Data       []T
	Pagination Pagination
}

// ListQuery represent the list request
type ListQuery struct {
	Page     int
	PageSize int
	Sort     string
	Query    string
}

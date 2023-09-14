package view

// Metadata is the response for metadata
type Metadata struct {
	Page         int    `json:"page" validate:"required"`
	PageSize     int    `json:"pageSize" validate:"required"`
	TotalPages   int    `json:"totalPages" validate:"required"`
	TotalRecords int    `json:"totalRecords" validate:"required"`
	Sort         string `json:"sort,omitempty"`
	HasNext      bool   `json:"hasNext,omitempty"`
} // @name Metadata

package util

import (
	"strings"

	"github.com/dwarvesf/go-api/pkg/model"
)

const (
	maxPageSize     = 1000
	defaultPageSize = 10
)

// CalculatePagination calculate pagination
func CalculatePagination(totalRecords int, page int, pageSize int) (*model.Pagination, error) {
	if page <= 0 {
		page = 1
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	hasNextRecords := totalRecords > page*pageSize

	return &model.Pagination{
		PageSize:     pageSize,
		Page:         page,
		TotalRecords: totalRecords,
		TotalPages:   calculateTotalPages(totalRecords, pageSize),
		Offset:       calculateOffset(page, pageSize),
		HasNext:      hasNextRecords,
	}, nil
}

func calculateOffset(page int, pageSize int) int {
	if page <= 0 {
		page = 1
	}

	return (page - 1) * pageSize
}

func calculateTotalPages(totalRecords int, pageSize int) int {
	// totalPages is grown up from int(count) / pageSize.
	// divide 0 will cause panic, so we need to check if pageSize is 0
	if pageSize == 0 {
		return 0
	}
	// if int(count) % pageSize != 0, totalPages will be increased by 1
	totalPages := totalRecords / pageSize
	if totalRecords%pageSize != 0 {
		totalPages++
	}

	return totalPages
}

// ParseSort parse sort string
func ParseSort(sort string) string {
	if sort == "" {
		return "created_at desc"
	}
	sItems := strings.Split(sort, ",")
	for i, s := range sItems {
		// remove space
		itm := strings.TrimSpace(s)
		// lower case
		itm = strings.ToLower(itm)
		if strings.HasPrefix(itm, "-") {
			sItems[i] = strings.TrimPrefix(itm, "-") + " desc"
		}

		if strings.HasPrefix(itm, "+") {
			sItems[i] = strings.TrimPrefix(itm, "+") + " asc"
		}
	}

	return strings.Join(sItems, ", ")
}

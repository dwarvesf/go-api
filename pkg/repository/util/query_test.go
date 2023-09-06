package util

import (
	"reflect"
	"testing"

	"github.com/dwarvesf/go-api/pkg/model"
)

func TestCalculatePagination(t *testing.T) {
	type args struct {
		totalRecords int
		page         int
		pageSize     int
	}
	tests := map[string]struct {
		args    args
		want    *model.Pagination
		wantErr bool
	}{
		"response with page 1": {
			args: args{
				totalRecords: 8,
				page:         1,
				pageSize:     10,
			},
			want: &model.Pagination{
				PageSize:     10,
				Page:         1,
				TotalRecords: 8,
				TotalPages:   1,
				Offset:       0,
			},
			wantErr: false,
		},
		"response with page 2 and pageSize 5": {
			args: args{
				totalRecords: 8,
				page:         2,
				pageSize:     5,
			},
			want: &model.Pagination{
				PageSize:     5,
				Page:         2,
				TotalRecords: 8,
				TotalPages:   2,
				Offset:       5,
			},
		},
		"page 0": {
			args: args{
				totalRecords: 8,
				page:         0,
				pageSize:     10,
			},
			want: &model.Pagination{
				PageSize:     10,
				Page:         1,
				TotalRecords: 8,
				TotalPages:   1,
				Offset:       0,
			},
		},
		"pageSize 999999": {
			args: args{
				totalRecords: 8,
				page:         1,
				pageSize:     999999,
			},
			want: &model.Pagination{
				PageSize:     1000,
				Page:         1,
				TotalRecords: 8,
				TotalPages:   1,
				Offset:       0,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := CalculatePagination(tt.args.totalRecords, tt.args.page, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("CaculatePagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CaculatePagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateTotalPages(t *testing.T) {
	type args struct {
		totalRecords int
		pageSize     int
	}
	tests := map[string]struct {
		name string
		args args
		want int
	}{
		"odd totalRecords": {
			args: args{
				totalRecords: 10,
				pageSize:     5,
			},
			want: 2,
		},
		"even totalRecords": {
			args: args{
				totalRecords: 11,
				pageSize:     2,
			},
			want: 6,
		},
		"zero totalRecords": {
			args: args{
				totalRecords: 0,
				pageSize:     2,
			},
			want: 0,
		},
		"zero pageSize": {
			args: args{
				totalRecords: 10,
				pageSize:     0,
			},
			want: 0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := calculateTotalPages(tt.args.totalRecords, tt.args.pageSize); got != tt.want {
				t.Errorf("calculateTotalPages() = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_calculateOffset(t *testing.T) {
	type args struct {
		page     int
		pageSize int
	}
	tests := map[string]struct {
		name string
		args args
		want int
	}{
		"page 1": {
			args: args{
				page:     1,
				pageSize: 10,
			},
			want: 0,
		},
		"page 0": {
			args: args{
				page:     0,
				pageSize: 10,
			},
			want: 0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := calculateOffset(tt.args.page, tt.args.pageSize); got != tt.want {
				t.Errorf("calculateOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSort(t *testing.T) {
	type args struct {
		sort string
	}
	tests := map[string]struct {
		name string
		args args
		want string
	}{
		"empty sort": {
			args: args{
				sort: "",
			},
			want: "created_at desc",
		},
		"sort with +": {
			args: args{
				sort: "+created_at",
			},
			want: "created_at asc",
		},
		"sort with -": {
			args: args{
				sort: "-created_at",
			},
			want: "created_at desc",
		},
		"sort with multiple fields": {
			args: args{
				sort: "+created_at,-updated_at",
			},
			want: "created_at asc, updated_at desc",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := ParseSort(tt.args.sort); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

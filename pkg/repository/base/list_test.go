package base

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ormUser struct {
	ID   int
	Name string
}

type modelUser struct {
	ID   int
	Name string
}

type counable struct {
	Value int64
	Error error
}

func (c counable) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if c.Error != nil {
		return 0, c.Error
	}
	return c.Value, nil
}

func mapping(o *ormUser) *modelUser {
	return &modelUser{
		ID:   o.ID,
		Name: o.Name,
	}
}

func TestGetList(t *testing.T) {
	type args struct {
		ctx db.Context
		q   model.ListQuery
		fns GetListFuncSet[ormUser, modelUser]
	}
	tests := map[string]struct {
		args    args
		want    *model.ListResult[modelUser]
		wantErr bool
	}{
		"contains data": {
			args: args{
				ctx: db.Context{},
				q: model.ListQuery{
					Page:     1,
					PageSize: 10,
					Sort:     "name",
					Query:    "test",
				},
				fns: GetListFuncSet[ormUser, modelUser]{
					PrepareQueryFn: func(ctx db.Context, q model.ListQuery) []qm.QueryMod {
						return []qm.QueryMod{}
					},
					CounableFn: func(qm []qm.QueryMod) Counable {
						return counable{Value: 1}
					},
					QueryListFn: func(qm []qm.QueryMod) ([]*ormUser, error) {
						return []*ormUser{
							{
								ID:   1,
								Name: "test",
							},
						}, nil
					},
					MappingFn: mapping,
				},
			},
			want: &model.ListResult[modelUser]{
				Data: []modelUser{
					{
						ID:   1,
						Name: "test",
					},
				},
				Pagination: model.Pagination{
					PageSize:     10,
					Page:         1,
					Sort:         "name",
					TotalRecords: 1,
					TotalPages:   1,
					Offset:       0,
					HasNext:      false,
				},
			},
			wantErr: false,
		},
		"empty data": {
			args: args{
				ctx: db.Context{},
				q: model.ListQuery{
					Page:     1,
					PageSize: 10,
					Sort:     "name",
					Query:    "test",
				},
				fns: GetListFuncSet[ormUser, modelUser]{
					PrepareQueryFn: func(ctx db.Context, q model.ListQuery) []qm.QueryMod {
						return []qm.QueryMod{}
					},
					CounableFn: func(qm []qm.QueryMod) Counable {
						return counable{Value: 0}
					},
					QueryListFn: func(qm []qm.QueryMod) ([]*ormUser, error) {
						return []*ormUser{}, nil
					},
					MappingFn: mapping,
				},
			},
			want: &model.ListResult[modelUser]{
				Pagination: model.Pagination{
					PageSize: 10,
					Page:     1,
					Sort:     "name",
				},
			},
			wantErr: false,
		},
		"query list error": {
			args: args{
				ctx: db.Context{},
				q: model.ListQuery{
					Page:     1,
					PageSize: 10,
					Sort:     "name",
					Query:    "test",
				},
				fns: GetListFuncSet[ormUser, modelUser]{
					PrepareQueryFn: func(ctx db.Context, q model.ListQuery) []qm.QueryMod {
						return []qm.QueryMod{}
					},
					CounableFn: func(qm []qm.QueryMod) Counable {
						return counable{Value: 0}
					},
					QueryListFn: func(qm []qm.QueryMod) ([]*ormUser, error) {
						return nil, errors.New("error")
					},
					MappingFn: mapping,
				},
			},
			want:    nil,
			wantErr: true,
		},
		"counable error": {
			args: args{
				ctx: db.Context{},
				q: model.ListQuery{
					Page:     1,
					PageSize: 10,
					Sort:     "name",
					Query:    "test",
				},
				fns: GetListFuncSet[ormUser, modelUser]{
					PrepareQueryFn: func(ctx db.Context, q model.ListQuery) []qm.QueryMod {
						return []qm.QueryMod{}
					},
					CounableFn: func(qm []qm.QueryMod) Counable {
						return counable{Value: 0, Error: errors.New("error")}
					},
					QueryListFn: func(qm []qm.QueryMod) ([]*ormUser, error) {
						return []*ormUser{}, nil
					},
					MappingFn: mapping,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetList(tt.args.ctx, tt.args.q, tt.args.fns)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}

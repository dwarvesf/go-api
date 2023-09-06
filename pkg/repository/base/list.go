package base

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/dwarvesf/go-api/pkg/repository/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Counable represent the countable interface
type Counable interface {
	Count(ctx context.Context, exec boil.ContextExecutor) (int64, error)
}

// GetListFuncSet represent the get list processor
type GetListFuncSet[OrmModel any, Model any] struct {
	PrepareQueryFn func(ctx db.Context, q model.ListQuery) []qm.QueryMod
	CounableFn     func([]qm.QueryMod) Counable
	QueryListFn    func([]qm.QueryMod) ([]*OrmModel, error)
	MappingFn      func(o *OrmModel) *Model
}

// GetList return the list of model
func GetList[OrmModel any, Model any](
	ctx db.Context,
	q model.ListQuery,
	fns GetListFuncSet[OrmModel, Model]) (*model.ListResult[Model], error) {
	// we will use the function set to process the list query
	// 1. prepare query
	// 2. calculate pagination
	// 3. query list
	// 4. mapping the result to the response model

	// prepare query and calculate pagination
	ormParams := fns.PrepareQueryFn(ctx, q)
	var count int64
	c := fns.CounableFn(ormParams)
	count, err := c.Count(ctx.Context, ctx.DB)
	if err != nil {
		return nil, err
	}
	pagination, err := util.CalculatePagination(int(count), q.Page, q.PageSize)
	if err != nil {
		return nil, err
	}
	nomalizedSort := util.ParseSort(q.Sort)
	pagination.Sort = nomalizedSort

	// query list
	queryParams := make([]qm.QueryMod, 0, len(ormParams))
	copy(queryParams, ormParams)
	queryParams = append(queryParams,
		qm.OrderBy(nomalizedSort),
		qm.Limit(pagination.PageSize),
		qm.Offset(pagination.Offset),
	)

	dt, err := fns.QueryListFn(queryParams)
	if err != nil {
		return nil, err
	}

	// mapping to model
	var result []Model
	for _, d := range dt {
		itm := fns.MappingFn(d)
		if itm != nil {
			result = append(result, *itm)
		}
	}

	return &model.ListResult[Model]{
		Data:       result,
		Pagination: *pagination,
	}, nil
}

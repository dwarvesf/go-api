package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
)

func (c impl) UserList(ctx context.Context, req model.ListQuery) (*model.ListResult[model.User], error) {
	dbCtx := db.FromContext(ctx)
	rs, err := c.repo.User.GetList(dbCtx, req)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

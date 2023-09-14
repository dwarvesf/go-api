package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
)

func (c *impl) UpdateUser(ctx context.Context, user model.UpdateUserRequest) (*model.User, error) {
	if c.monitor != nil {
		const spanName = "UpdateUserController"
		newCtx, span := c.monitor.Start(ctx, spanName)
		ctx = newCtx
		defer span.End()
	}

	uID, err := middleware.UserIDFromContext(ctx)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	dbCtx := db.FromContext(ctx)
	u, err := c.repo.User.GetByID(dbCtx, uID)
	if err != nil {
		return nil, err
	}

	u.Avatar = user.Avatar
	u.FullName = user.FullName

	updated, err := c.repo.User.Update(dbCtx, uID, user)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

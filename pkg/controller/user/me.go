package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
)

func (c *impl) Me(ctx context.Context) (*model.User, error) {
	if c.monitor != nil {
		const spanName = "MeController"
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

	return u, nil
}

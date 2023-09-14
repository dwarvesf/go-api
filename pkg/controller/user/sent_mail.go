package user

import (
	"context"
	"fmt"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
)

const pageSize = 10

func (c *impl) SentMail(ctx context.Context) error {
	if c.monitor != nil {
		const spanName = "LoginController"
		newCtx, span := c.monitor.Start(ctx, spanName)
		ctx = newCtx
		defer span.End()
	}

	dbCtx := db.FromContext(ctx)

	hashNext := true
	page := 1
	// while has next page
	for hashNext {
		userList, err := c.repo.User.GetList(dbCtx,
			model.ListQuery{
				Page:     page,
				PageSize: pageSize,
			},
		)
		if err != nil {
			return err
		}

		for _, user := range userList.Data {
			// fake sent mail to user
			fmt.Println("Sent mail to user: ", user.Email)
		}

		hashNext = userList.Pagination.HasNext
		page++
	}

	return nil
}

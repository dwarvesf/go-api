package user

import (
	"context"
	"fmt"

	"github.com/dwarvesf/go-api/pkg/repository/db"
)

const pageSize = 10

func (c *impl) SentMail(ctx context.Context) error {
	dbCtx := db.FromContext(ctx)

	hashNext := true
	page := 1
	// while has next page
	for hashNext {
		userList, err := c.repo.User.GetList(dbCtx, page, pageSize, "", "")
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

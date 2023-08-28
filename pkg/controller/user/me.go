package user

import (
	"github.com/dwarvesf/go-api/pkg/repository/orm"
)

func (c *impl) Me(userID int) (*orm.User, error) {
	// TODO: get user from database
	var u *orm.User

	return u, nil
}

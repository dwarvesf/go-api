package user

import (
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
)

func (c *impl) UpdateUser(ID int, user model.UpdateUserRequest) (*orm.User, error) {
	// TODO: update user to database
	var u *orm.User

	return u, nil
}

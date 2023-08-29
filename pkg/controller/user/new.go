package user

import (
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
)

// Controller auth controller
type Controller interface {
	Me(userID int) (*orm.User, error)
}

type impl struct {
	repo orm.Repo
	cfg  config.Config
}

// NewUserController new auth controller
func NewUserController(cfg config.Config, r orm.Repo) Controller {
	return &impl{
		repo: r,
		cfg:  cfg,
	}
}

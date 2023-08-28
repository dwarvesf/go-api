package user

import (
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
)

// Controller auth controller
//
//go:generate mockery --name=Controller --with-expecter --output ./mocks
type Controller interface {
	Me(userID int) (*orm.User, error)
}

type impl struct {
	repo orm.Repo
	cfg  config.Config
}

func NewUserController(cfg config.Config, r orm.Repo) Controller {
	return &impl{
		repo: r,
		cfg:  cfg,
	}
}

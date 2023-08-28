package auth

import (
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
	"github.com/dwarvesf/go-api/pkg/service/jwthelper"
)

// Controller auth controller
//
//go:generate mockery --name=Controller --with-expecter --output ./mocks
type Controller interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
	Signup(req model.SignupRequest) error
}

type impl struct {
	repo      orm.Repo
	jwtHelper jwthelper.Helper
	cfg       config.Config
}

func NewAuthController(cfg config.Config, r orm.Repo) Controller {
	return &impl{
		repo:      r,
		jwtHelper: jwthelper.NewHelper(cfg.SecretKey),
		cfg:       cfg,
	}
}

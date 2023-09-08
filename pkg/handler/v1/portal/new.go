package portal

import (
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/controller/auth"
	"github.com/dwarvesf/go-api/pkg/controller/user"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/logger/monitor"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/service"
)

// Handler for app
type Handler struct {
	cfg      config.Config
	log      logger.Log
	svc      service.Service
	monitor  monitor.Exporter
	authCtrl auth.Controller
	userCtrl user.Controller
}

// New will return an instance of Auth struct
func New(cfg config.Config, l logger.Log, repo *repository.Repo, svc service.Service, monitor monitor.Exporter) *Handler {
	return &Handler{
		cfg:      cfg,
		log:      l,
		svc:      svc,
		monitor:  monitor,
		authCtrl: auth.NewAuthController(cfg, repo, monitor),
		userCtrl: user.NewUserController(cfg, repo, monitor),
	}
}

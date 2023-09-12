package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/logger/monitor"
	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/realtime"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/dwarvesf/go-api/pkg/service"
	"github.com/dwarvesf/go-api/pkg/service/jwthelper"
)

// @title           APP API DOCUMENT
// @version         v0.0.1
// @description     This is api document for APP API project.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Andy
// @contact.url    https://d.foundation
// @contact.email  andy@d.foundation

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig(config.DefaultConfigLoaders())
	sMonitor, err := monitor.NewSentry(cfg)
	if err != nil {
		log.Fatal(err, "failed to init sentry")
	}
	defer sMonitor.Clean(2 * time.Second)

	l := logger.NewLogByConfig(cfg)
	l.Infof("Server starting")

	authMw := middleware.NewAuthMiddleware(jwthelper.NewHelper(cfg.SecretKey))
	a := App{
		l:              l,
		cfg:            cfg,
		service:        service.New(cfg),
		repo:           repository.NewRepo(),
		monitor:        sMonitor,
		realtimeServer: realtime.New(authMw, l),
	}

	_, err = db.Init(*cfg)
	if err != nil {
		l.Fatal(err, "failed to init db")
	}

	// Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: setupRouter(a),
	}

	quit := make(chan os.Signal)

	// serve http server
	go func() {
		a.l.Info("listening on " + a.cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err, "failed to listen and serve")
		}

		quit <- os.Interrupt
	}()

	signal.Notify(quit, os.Interrupt)

	<-quit

	shutdownServer(srv, l)
}

func shutdownServer(srv *http.Server, l logger.Log) {
	l.Info("Server Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		l.Error(err, "failed to shutdown server")
	}

	l.Info("Server Exit")
}

// App api app instance
type App struct {
	l              logger.Log
	cfg            *config.Config
	service        service.Service
	repo           *repository.Repo
	monitor        monitor.Tracer
	realtimeServer realtime.Server
}

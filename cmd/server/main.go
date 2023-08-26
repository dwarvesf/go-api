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
	"github.com/dwarvesf/go-api/pkg/service"
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

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig(config.DefaultConfigLoaders())
	sClient, err := logger.NewSentry(cfg)
	if err != nil {
		log.Fatal(err, "failed to init sentry")
	}
	if sClient != nil {
		defer sClient.Flush(2 * time.Second)
	}

	l := logger.NewLogByConfig(cfg, sClient)
	l.Infof("Server starting")

	a := App{
		l:   l,
		cfg: cfg,
	}

	// Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: setupRouter(a, sClient),
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
	l       logger.Log
	cfg     *config.Config
	service service.Service
}

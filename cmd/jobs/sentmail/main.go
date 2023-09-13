package main

import (
	"context"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/controller/user"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/logger/monitor"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
)

func main() {
	cfg := config.LoadConfig(config.DefaultConfigLoaders())

	l := logger.NewLogByConfig(cfg)
	l.Infof("Cronjob starting")

	_, err := db.Init(*cfg)
	if err != nil {
		l.Fatal(err, "failed to init db")
	}

	sentryMonitor, err := monitor.NewSentry(cfg)
	if err != nil {
		l.Fatal(err, "failed to init sentry")
	}
	defer sentryMonitor.Clean(2 * time.Second)

	// new controler
	c := user.NewUserController(*cfg, repository.NewRepo(), sentryMonitor)
	c.SentMail(context.Background())
}

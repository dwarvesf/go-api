package main

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/controller/user"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
)

func main() {
	cfg := config.LoadConfig(config.DefaultConfigLoaders())

	l := logger.NewLogByConfig(cfg, nil)
	l.Infof("Cronjob starting")

	_, err := db.Init(*cfg)
	if err != nil {
		l.Fatal(err, "failed to init db")
	}

	// new controler
	c := user.NewUserController(*cfg, repository.NewRepo())
	c.SentMail(context.Background())
}

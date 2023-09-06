package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// WithTestingDB run callback with transaction
func WithTestingDB(t *testing.T, callback func(ctx Context)) {
	if singleDB == nil {
		initTestingDB(t)
	}

	tx, err := singleDB.Begin()
	require.NoError(t, err)

	defer tx.Rollback()

	callback(Context{
		Context: context.Background(),
		DB:      tx,
	})
}

func initTestingDB(t *testing.T) {
	cfg := config.LoadTestConfig()
	dbConnCfg, err := pgx.ParseConfig(cfg.DatabaseURL)
	require.NoError(t, err)
	connStr := stdlib.RegisterConnConfig(dbConnCfg)
	appDB, err := sql.Open("pgx", connStr)
	appDB.SetMaxOpenConns(50)
	appDB.SetConnMaxLifetime(30 * time.Minute)

	setupGlobalDB(appDB)
}

package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

// singleDB is global db
var singleDB *sql.DB

// once is used for setup global db
var once sync.Once

// setupGlobalDB setup global db
func setupGlobalDB(db *sql.DB) {
	once.Do(func() {
		singleDB = db
	})
}

// GetDB get global db
func GetDB() *sql.DB {
	if singleDB == nil {
		panic("db is not initialized")
	}

	return singleDB
}

// Init create new db pool
func Init(cfg config.Config) (*sql.DB, error) {
	connCfg, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("parsing pgx config failed. err: %w", err))
	}

	pool, err := sql.Open("pgx", stdlib.RegisterConnConfig(connCfg))
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("opening DB failed. err: %w", err))
	}

	pool.SetConnMaxLifetime(29 * time.Minute)
	pool.SetMaxOpenConns(cfg.DBMaxOpenConns)
	pool.SetMaxIdleConns(cfg.DBMaxIdleConns)

	setupGlobalDB(pool)

	return pool, nil
}

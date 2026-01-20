package db

import (
	"context"
	"database/sql"
	"time"

	appcfg "mutabaahapi/internal/infra/config"
	"mutabaahapi/internal/platform/dbx"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// DB wraps *sqlx.DB and provides a simple transaction helper.
type DB struct {
	DB *sqlx.DB
}

func New(cfg appcfg.Config) (*DB, error) {
	database, err := sqlx.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Pool tuning from config
	database.SetMaxOpenConns(cfg.DBMaxOpenConns)
	database.SetMaxIdleConns(cfg.DBMaxIdleConns)
	database.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetimeMin) * time.Minute)
	database.SetConnMaxIdleTime(time.Duration(cfg.DBConnMaxIdleTimeMin) * time.Minute)

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, err
	}
	return &DB{DB: database}, nil
}

func (d *DB) Close() error { return d.DB.Close() }

// WithinTx runs fn inside a transaction with default options and provides a dbx.Executor.
func (d *DB) WithinTx(ctx context.Context, fn func(ctx context.Context, q dbx.Executor) error) error {
	tx, err := d.DB.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

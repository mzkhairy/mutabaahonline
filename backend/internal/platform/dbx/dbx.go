package dbx

import (
    "context"
    "database/sql"

    "github.com/jmoiron/sqlx"
)

// Executor is a minimal interface implemented by *sqlx.DB and *sqlx.Tx.
type Executor interface {
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
    QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
    GetContext(ctx context.Context, dest any, query string, args ...any) error
    SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

// UnitOfWork runs a function within a transaction and provides an Executor bound to it.
type UnitOfWork interface {
    WithinTx(ctx context.Context, fn func(ctx context.Context, q Executor) error) error
}


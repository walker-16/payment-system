package db

import (
	"context"
	"time"

	"github.com/jackc/pgx"
)

// Config holds database connection settings and pool configuration.
type Config struct {
	DSN             string        `env:"DB_DSN,required"`
	MaxConns        int32         `env:"DB_MAX_CONNS,default=10"`
	MinConns        int32         `env:"DB_MIN_CONNS,default=0"`
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME,default=30m"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME,default=1h"`
	AppName         string        `env:"DB_APP_NAME"`
}

// DB defines the interface for database operations.
type DB interface {
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	QueryRow(ctx context.Context, dest any, query string, args ...any) error
	BeginTx(ctx context.Context) (Tx, error)
	Ping(ctx context.Context) error
	Close()
}

// Tx defines the interface for transactional operations.
type Tx interface {
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Queryable defines low-level database operations compatible with pgx.
type Queryable interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgx.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Close()
}

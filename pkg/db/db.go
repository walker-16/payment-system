package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	pool *pgxpool.Pool
}

type TxOptions struct{}

// New initializes a new database connection pool using the provided configuration.
func New(ctx context.Context, cfg Config) (DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("DSN cannot be empty")
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database url: %w", err)
	}
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.ConnConfig.RuntimeParams["application_name"] = cfg.AppName

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &db{pool: pool}, nil
}

// Close closes the database connection pool.
func (c *db) Close() {
	c.pool.Close()
}

// Ping verifies a connection to the database is still alive.
func (c *db) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

// Select queries the database and maps the results to a slice of structs.
func (c *db) Select(ctx context.Context, dest any, query string, args ...any) error {
	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	if err := pgxscan.ScanAll(dest, rows); err != nil {
		return fmt.Errorf("failed to scan all rows: %w", err)
	}
	return nil
}

// Exec executes a query without returning any rows.
func (c *db) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	commandTag, err := c.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	return commandTag.RowsAffected(), nil
}

// QueryRow executes a query and scans the result into the provided struct.
func (c *db) QueryRow(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Get(ctx, c.pool, dest, query, args...)
}

// BeginTx starts a new transaction and returns a Tx interface.
func (c *db) BeginTx(ctx context.Context) (Tx, error) {
	txr, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &tx{tx: txr}, nil
}

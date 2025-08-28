package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Tx interface {
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type tx struct {
	tx pgx.Tx
}

func (t *tx) Select(ctx context.Context, dest any, query string, args ...any) error {
	rows, err := t.tx.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute transaction query: %w", err)
	}
	defer rows.Close()

	if err := pgxscan.ScanAll(dest, rows); err != nil {
		return fmt.Errorf("failed to scan transaction rows: %w", err)
	}
	return nil
}

func (t *tx) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	commandTag, err := t.tx.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute transaction query: %w", err)
	}
	return commandTag.RowsAffected(), nil
}

func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

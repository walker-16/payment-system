package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type tx struct {
	tx pgx.Tx
}

// Select executes a query within the transaction and scans all resulting rows
// into the provided destination slice. Returns an error if the query fails
// or scanning fails.
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

// Exec executes a query within the transaction without returning rows.
// Returns the number of affected rows or an error if execution fails.
func (t *tx) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	commandTag, err := t.tx.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute transaction query: %w", err)
	}
	return commandTag.RowsAffected(), nil
}

// Commit commits the current transaction. Returns an error if commit fails.
func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

// Rollback aborts the current transaction and rolls back any changes.
// Returns an error if rollback fails.
func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	ports "github.com/koubae/game-hangar/pkg/database"
)

type PostgresTransaction struct {
	tx txInterface
}

var _ ports.Transaction = (*PostgresTransaction)(nil)

func (t *PostgresTransaction) SelectMany(
	ctx context.Context,
	query string,
	args ...any,
) (pgx.Rows, error) {
	return t.tx.Query(ctx, query, args...)
}

func (t *PostgresTransaction) SelectOne(
	ctx context.Context,
	query string,
	args ...any,
) pgx.Row {
	return t.tx.QueryRow(ctx, query, args...)
}

func (t *PostgresTransaction) SQL(
	ctx context.Context,
	query string,
	args ...any,
) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, query, args...)
}

func (t *PostgresTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *PostgresTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

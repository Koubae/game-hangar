package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/koubae/game-hangar/pkg/database"
)

func MapPostgresErrToDomainErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return database.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505": // unique_violation
		return &database.ErrDuplicate{
			Err:        err,
			Constraint: pgErr.ConstraintName,
		}
	default:
		return err
	}
}

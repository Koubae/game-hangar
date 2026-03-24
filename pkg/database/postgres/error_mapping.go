package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/koubae/game-hangar/pkg/database"
)

func MapPostgresErrToDomainErr(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505": // unique_violation
		return fmt.Errorf("unique violation %w (%s): %w", database.ErrrDuplicate, pgErr.ConstraintName, err)
	default:
		return err
	}
}

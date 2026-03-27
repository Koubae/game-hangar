package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("resource not found")

type ErrDuplicate struct {
	Err        error // underlying DB error
	Constraint string
}

func (e *ErrDuplicate) Error() string {
	return fmt.Sprintf("unique violation on %s: %v", e.Constraint, e.Err)
}

func (e *ErrDuplicate) Unwrap() error {
	return e.Err
}

func (e *ErrDuplicate) Is(target error) bool { // NOTE: Sentinel struct pattern
	_, ok := target.(*ErrDuplicate)
	return ok
}

type ErrOpenTransaction struct {
	Err error // underlying DB error
}

func (e *ErrOpenTransaction) Error() string {
	return fmt.Sprintf("error opening a new DB transaction, error: %v", e.Err)
}

func (e *ErrOpenTransaction) Unwrap() error {
	return e.Err
}

func (e *ErrOpenTransaction) Is(target error) bool {
	_, ok := target.(*ErrOpenTransaction)
	return ok
}

// DBTX pool-backed connector and a transaction wrapper implement this.
type DBTX interface {
	SelectMany(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	SelectOne(ctx context.Context, query string, args ...any) pgx.Row
	SQL(
		ctx context.Context,
		query string,
		args ...any,
	) (pgconn.CommandTag, error)
	MapDBErrToDomainErr(err error) error
}

type Transaction interface {
	DBTX
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Connector is the app-facing database entrypoint.
type Connector interface {
	DBTX

	String() string
	Ping(ctx context.Context) error
	Shutdown() error
	Transaction(
		ctx context.Context,
		txOptions pgx.TxOptions,
	) (Transaction, error)
}

package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound   = errors.New("resource not found")
	ErrrDuplicate = errors.New("resource is duplicate")
)

// DBTX pool-backed connector and a transaction wrapper implement this.
type DBTX interface {
	SelectMany(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	SelectOne(ctx context.Context, query string, args ...any) pgx.Row
	SQL(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
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
	Transaction(ctx context.Context, txOptions pgx.TxOptions) (Transaction, error)
}

package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/stretchr/testify/mock"
)

var (
	Now             = time.Now()
	AccountIDTest01 = uuid.New()

	DBMockErrDuplicateKey = &pgconn.PgError{
		Code:           "23505",
		ConstraintName: "some_unique_constraint_name",
		Message:        "duplicate key value violates unique constraint",
	}
)

type MockDBPool struct {
	mock.Mock
}

var _ postgres.PoolInterface = (*MockDBPool)(nil)

func (m *MockDBPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDBPool) Close() {
	m.Called()
}

func (m *MockDBPool) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	callArgs := m.Called(append([]any{ctx, query}, args...)...)

	var rows pgx.Rows
	if v := callArgs.Get(0); v != nil {
		rows = v.(pgx.Rows)
	}

	return rows, callArgs.Error(1)
}

func (m *MockDBPool) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	callArgs := m.Called(append([]any{ctx, query}, args...)...)

	if v := callArgs.Get(0); v != nil {
		return v.(pgx.Row)
	}

	return nil
}

func (m *MockDBPool) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	callArgs := m.Called(append([]any{ctx, query}, args...)...)

	var tag pgconn.CommandTag
	if v := callArgs.Get(0); v != nil {
		tag = v.(pgconn.CommandTag)
	}

	return tag, callArgs.Error(1)
}

func (m *MockDBPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	callArgs := m.Called(ctx, txOptions)

	var tx pgx.Tx
	if v := callArgs.Get(0); v != nil {
		tx = v.(pgx.Tx)
	}

	return tx, callArgs.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Args(n int) []any {
	args := make([]any, n)
	for i := range n {
		args[i] = mock.Anything
	}
	return args
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRow) MockScan(argsN int, err error, values ...any) {
	m.On("Scan", m.Args(argsN)...).Run(func(args mock.Arguments) {
		if err != nil {
			return
		}

		set := func(ptr any, _index, val any) {
			switch ptr := ptr.(type) {
			case *int:
				*ptr = val.(int)
			case **int:
				*ptr = val.(*int)

			case *int64:
				*ptr = val.(int64)
			case **int64:
				*ptr = val.(*int64)

			case *string:
				*ptr = val.(string)
			case **string:
				*ptr = val.(*string)

			case *bool:
				*ptr = val.(bool)
			case **bool:
				*ptr = val.(*bool)

			case *uuid.UUID:
				*ptr = val.(uuid.UUID)
			case **uuid.UUID:
				*ptr = val.(*uuid.UUID)

			case *time.Time:
				*ptr = val.(time.Time)
			case **time.Time:
				*ptr = val.(*time.Time)

			case *any:
				*ptr = val

			default:
				panic(fmt.Sprintf("MockScan: untyped or unhandled destination at index %d (%T)", _index, ptr))
			}
		}

		setNil := func(ptr any) {
			switch p := ptr.(type) {
			case **string:
				*p = nil
			case **int64:
				*p = nil
			case **bool:
				*p = nil
			case **uuid.UUID:
				*p = nil
			case **time.Time:
				*p = nil
			}
		}

		for i, val := range values {
			if i >= len(args) {
				break
			}

			ptr := args.Get(i)
			if val == nil {
				setNil(ptr)
				continue
			}

			set(ptr, i, val)
		}
	}).Return(err)
}

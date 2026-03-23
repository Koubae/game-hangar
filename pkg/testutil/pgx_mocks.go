package testutil

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/stretchr/testify/mock"
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
		set := func(_index int, val any) {
			switch ptr := args.Get(_index).(type) {
			case *int:
				*ptr = val.(int)
			case *string:
				*ptr = val.(string)
			case *bool:
				*ptr = val.(bool)
			}
		}

		for i, val := range values {
			set(i, val)
		}
	}).Return(err)
}

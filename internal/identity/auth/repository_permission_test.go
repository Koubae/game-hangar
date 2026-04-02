package auth_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPermissionRepository_GetPermissions_CacheHit(t *testing.T) {
	t.Parallel()
	common.CreateLogger(common.LogLevelDPanic, "")

	expected := []*auth.Permission{
		{
			ID:       1,
			Service:  "identity",
			Resource: "account",
			Action:   "read",
			Created:  testutil.Now,
		},
		{
			ID:       2,
			Service:  "identity",
			Resource: "account",
			Action:   "write",
			Created:  testutil.Now,
		},
		{
			ID:       3,
			Service:  "identity",
			Resource: "account",
			Action:   "delete",
			Created:  testutil.Now,
		},
	}

	mockPool := new(testutil.MockDBPool)
	connector := postgres.ConnectorPostgres{Pool: mockPool}

	repo := &auth.PermissionRepository{
		PermissionsCache: map[int64]*auth.Permission{
			1: expected[0],
			2: expected[1],
			3: expected[2],
		},
	}
	got := repo.GetPermissions(
		context.Background(),
		&connector,
		[]int64{1, 2, 3},
	)

	assert.Equal(t, expected, got)
	mockPool.AssertExpectations(t)
}

func TestPermissionRepository_GetPermissions_CacheMiss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id       string
		expected []*auth.Permission
	}{
		{
			id: "record-is-found",
			expected: []*auth.Permission{
				{
					ID:       1,
					Service:  "identity",
					Resource: "account",
					Action:   "read",
					Created:  testutil.Now,
				},
			},
		},

		{
			id:       "record-is-not-found",
			expected: []*auth.Permission{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				common.CreateLogger(common.LogLevelDPanic, "")
				mockRow := new(testutil.MockRow)

				if len(tt.expected) > 0 {
					mockRow.MockScan(
						5, nil,
						tt.expected[0].ID,
						tt.expected[0].Service,
						tt.expected[0].Resource,
						tt.expected[0].Action,
						tt.expected[0].Created,
					)

				} else {
					mockRow.MockScan(5, pgx.ErrNoRows)
				}

				mockPool := new(testutil.MockDBPool)
				mockPool.On(
					"QueryRow", mock.Anything, mock.Anything, pgx.StrictNamedArgs{
						"id": int64(1),
					},
				).
					Return(mockRow)

				connector := postgres.ConnectorPostgres{Pool: mockPool}

				repo := &auth.PermissionRepository{
					PermissionsCache: map[int64]*auth.Permission{},
				}

				got := repo.GetPermissions(
					context.Background(),
					&connector,
					[]int64{1},
				)

				assert.Equal(t, tt.expected, got)
				mockPool.AssertExpectations(t)
			},
		)
	}
}

func TestPermissionRepository_GetAdminAccountPermissions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		accountID  string
		rows       [][]any
		errOnQuery error
		errOnScan  error
		expected   []*auth.Permission
	}{
		"records-are-found": {
			accountID:  testunit.AccountIDTest01Str,
			errOnQuery: nil,
			errOnScan:  nil,
			expected: []*auth.Permission{
				{
					ID:       1,
					Service:  "identity",
					Resource: "account",
					Action:   "read",
					Created:  testutil.Now,
				},
				{
					ID:       2,
					Service:  "identity",
					Resource: "account",
					Action:   "write",
					Created:  testutil.Now,
				},
				{
					ID:       3,
					Service:  "identity",
					Resource: "account",
					Action:   "delete",
					Created:  testutil.Now,
				},
				{
					ID:       4,
					Service:  "identity",
					Resource: "account",
					Action:   "*",
					Created:  testutil.Now,
				},
			},

			rows: [][]any{
				{int64(1), "identity", "account", "read", testutil.Now},
				{int64(2), "identity", "account", "write", testutil.Now},
				{int64(3), "identity", "account", "delete", testutil.Now},
				{int64(4), "identity", "account", "*", testutil.Now},
			},
		},
		"records-are-not-found": {
			accountID:  testunit.AccountIDTest01Str,
			errOnQuery: nil,
			errOnScan:  nil,
			expected:   []*auth.Permission{},
			rows:       [][]any{},
		},
		"on-err-query": {
			accountID:  testunit.AccountIDTest01Str,
			errOnQuery: testunit.ErrDBGeneric,
			errOnScan:  nil,
			expected:   []*auth.Permission{},
			rows: [][]any{
				{int64(1), "identity", "account", "read", testutil.Now},
				{int64(2), "identity", "account", "write", testutil.Now},
				{int64(3), "identity", "account", "delete", testutil.Now},
				{int64(4), "identity", "account", "*", testutil.Now},
			},
		},
		"on-err-scan": {
			accountID:  testunit.AccountIDTest01Str,
			errOnQuery: nil,
			errOnScan:  testunit.ErrDBGeneric,
			expected:   []*auth.Permission{},
			rows: [][]any{
				{int64(1), "identity", "account", "read", testutil.Now},
				{int64(2), "identity", "account", "write", testutil.Now},
				{int64(3), "identity", "account", "delete", testutil.Now},
				{int64(4), "identity", "account", "*", testutil.Now},
			},
		},
	}

	common.CreateLogger(common.LogLevelDPanic, "")
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				mockRows := &testutil.MockRows{
					Data:    tt.rows,
					ScanErr: tt.errOnScan,
				}
				mockPool := new(testutil.MockDBPool)
				mockPool.On(
					"Query", mock.Anything, mock.Anything, pgx.StrictNamedArgs{
						"account_id": tt.accountID,
					},
				).
					Return(mockRows, tt.errOnQuery)

				connector := postgres.ConnectorPostgres{Pool: mockPool}
				repo := auth.NewPermissionRepository()

				got := repo.GetAdminAccountPermissions(
					context.Background(),
					&connector,
					tt.accountID,
				)

				assert.Equal(t, tt.expected, got)
				mockPool.AssertExpectations(t)
			},
		)
	}
}

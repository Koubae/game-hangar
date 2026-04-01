package auth_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/auth"
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

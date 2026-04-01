package repository

import (
	"context"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/tests/integration"
	"github.com/stretchr/testify/assert"
)

func TestPermissionRepository_GetPermissions(t *testing.T) {
	_, connector, tearDown := integration.DBWithCleanup(t, true)
	defer tearDown()

	tests := map[string]struct {
		ids      []int64
		expected []*auth.Permission
	}{
		"permission-su": {
			ids: []int64{1},
			expected: []*auth.Permission{
				{
					ID:       1,
					Service:  "*",
					Resource: "*",
					Action:   "*",
				},
			},
		},
		"permission-su-with-suidentity": {
			ids: []int64{1, 2},
			expected: []*auth.Permission{
				{
					ID:       1,
					Service:  "*",
					Resource: "*",
					Action:   "*",
				},
				{
					ID:       2,
					Service:  "identity",
					Resource: "*",
					Action:   "*",
				},
			},
		},

		"permission-all-identity-auth": {
			ids: []int64{3, 4, 5, 6},
			expected: []*auth.Permission{
				{
					ID:       3,
					Service:  "identity",
					Resource: "auth",
					Action:   "*",
				},
				{
					ID:       4,
					Service:  "identity",
					Resource: "auth",
					Action:   "read",
				},
				{
					ID:       5,
					Service:  "identity",
					Resource: "auth",
					Action:   "write",
				},
				{
					ID:       6,
					Service:  "identity",
					Resource: "auth",
					Action:   "delete",
				},
			},
		},

		"permission-all-identity-account": {
			ids: []int64{7, 8, 9, 10},
			expected: []*auth.Permission{
				{
					ID:       7,
					Service:  "identity",
					Resource: "account",
					Action:   "*",
				},
				{
					ID:       8,
					Service:  "identity",
					Resource: "account",
					Action:   "read",
				},
				{
					ID:       9,
					Service:  "identity",
					Resource: "account",
					Action:   "write",
				},
				{
					ID:       10,
					Service:  "identity",
					Resource: "account",
					Action:   "delete",
				},
			},
		},
	}

	permissionRepository := auth.NewPermissionRepository()
	permissionRepository.LoadPermissions(context.Background(), connector)
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				permissions := permissionRepository.GetPermissions(context.Background(), connector, tt.ids)

				for i, permission := range permissions {
					assert.Equal(t, tt.expected[i].Service, permission.Service)
					assert.Equal(t, tt.expected[i].Resource, permission.Resource)
					assert.Equal(t, tt.expected[i].Action, permission.Action)
				}
			},
		)
	}
}

func TestPermissionRepository_GetPermissionsFoundWhenCacheMiss(t *testing.T) {
	_, connector, tearDown := integration.DBWithCleanup(t, true)
	defer tearDown()

	ids := []int64{3, 4, 5, 6}
	permissionRepository := auth.NewPermissionRepository()
	permissions := permissionRepository.GetPermissions(context.Background(), connector, ids)

	expected := []auth.Permission{
		{Service: "identity", Resource: "auth", Action: "*"},
		{Service: "identity", Resource: "auth", Action: "read"},
		{Service: "identity", Resource: "auth", Action: "write"},
		{Service: "identity", Resource: "auth", Action: "delete"},
	}

	for i, permission := range permissions {
		assert.Equal(t, expected[i].Service, permission.Service)
		assert.Equal(t, expected[i].Resource, permission.Resource)
		assert.Equal(t, expected[i].Action, permission.Action)
	}

}

func TestPermissionRepository_GetPermissionsNotFound(t *testing.T) {
	_, connector, tearDown := integration.DBWithCleanup(t, true)
	defer tearDown()

	ids := []int64{99, 100, 101, 102}
	permissionRepository := auth.NewPermissionRepository()
	permissions := permissionRepository.GetPermissions(context.Background(), connector, ids)

	expected := make([]*auth.Permission, 0)
	assert.Equal(t, expected, permissions)
}

package authpkg_test

import (
	"testing"

	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"github.com/stretchr/testify/assert"
)

func TestScope_ParseScope(t *testing.T) {
	tests := map[string]struct {
		scope       string
		expected    *authpkg.Scope
		expectedErr error
	}{
		"valid-scope-read": {
			scope: "identity:account:read",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"read"},
			},
			expectedErr: nil,
		},
		"valid-scope-write": {
			scope: "identity:account:write",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"write"},
			},
			expectedErr: nil,
		},
		"valid-scope-delete": {
			scope: "identity:account:delete",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"delete"},
			},
			expectedErr: nil,
		},
		"valid-scope-read-write": {
			scope: "identity:account:read,write",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"read", "write"},
			},
			expectedErr: nil,
		},
		"valid-scope-read-delete": {
			scope: "identity:account:read,delete",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"delete", "read"},
			},
			expectedErr: nil,
		},
		"valid-scope-write-delete": {
			scope: "identity:account:write,delete",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"delete", "write"},
			},
			expectedErr: nil,
		},
		"valid-scope-only-wildcard-if-present": {
			scope: "identity:account:read,write,*",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"*"},
			},
			expectedErr: nil,
		},
		"valid-scope-spaces-and-uppercase-normalized": {
			scope: "IdenTity : ACCOUNT : WrITe, DeletE, READ,   ",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"delete", "read", "write"},
			},
			expectedErr: nil,
		},

		"invalid-scope-action-is-wildcard": {
			scope: "identity:account:*",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Actions:  []authpkg.Action{"*"},
			},
			expectedErr: nil,
		},
		"invalid-scope-resource-is-wildcard": {
			scope: "identity:*:read",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "*",
				Actions:  []authpkg.Action{"read"},
			},
			expectedErr: nil,
		},
		"invalid-scope-service-is-wildcard": {
			scope: "*:account:read",
			expected: &authpkg.Scope{
				Service:  "*",
				Resource: "account",
				Actions:  []authpkg.Action{"read"},
			},
			expectedErr: nil,
		},

		"invalid-scope-action-cannot-exceed-more-than-1-empty-space": {
			scope:       "identity:account:read,write,delete,,",
			expected:    nil,
			expectedErr: errspkg.AuthPermissionsScopeFormat,
		},
		"invalid-scope-action-cannot-exceed-3-actions": {
			scope:       "identity:account:read,write,delete,read",
			expected:    nil,
			expectedErr: errspkg.AuthPermissionsScopeFormat,
		},

		"invalid-scope-action-is-one-char-and-not-a-wildcard": {
			scope:       "identity:account:!",
			expected:    nil,
			expectedErr: errspkg.AuthPermissionsScopeFormat,
		},
		"invalid-scope-resource-is-one-char-and-not-a-wildcard": {
			scope:       "identity:!:read",
			expected:    nil,
			expectedErr: errspkg.AuthPermissionsScopeFormat,
		},
		"invalid-scope-service-is-one-char-and-not-a-wildcard": {
			scope:       "!:account:read",
			expected:    nil,
			expectedErr: errspkg.AuthPermissionsScopeFormat,
		},

		"invalid-scope-action-is-invalid": {
			scope:       "identity:account:update",
			expected:    nil,
			expectedErr: errspkg.AuthPermissionsScopeFormat,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				scope, err := authpkg.ParseScope(tt.scope)

				if tt.expectedErr != nil {
					assert.Nil(t, scope)
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.expectedErr)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.expected, scope)

			},
		)
	}

}

func TestPermissions_NewPermissions(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		scopes   []authpkg.Scope
		expected authpkg.Permissions
	}{
		"permissions-created-normal": {
			scopes: []authpkg.Scope{
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.READ, authpkg.WRITE}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "storage", Actions: []authpkg.Action{authpkg.READ, authpkg.WRITE}},
			},
			expected: authpkg.Permissions{
				"identity": {
					"account":             {authpkg.READ, authpkg.WRITE},
					"account_credentials": {authpkg.READ},
				},
				"storage": {
					"config":  {authpkg.READ},
					"storage": {authpkg.READ, authpkg.WRITE},
				},
			},
		},
		"permissions-merges-duplicate-resources": {
			scopes: []authpkg.Scope{
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.READ, authpkg.WRITE}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.READ, authpkg.WRITE}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.DELETE}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.DELETE}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.DELETE}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.DELETE}},
			},
			expected: authpkg.Permissions{
				"identity": {
					"account":             {authpkg.DELETE, authpkg.READ, authpkg.WRITE},
					"account_credentials": {authpkg.READ, authpkg.WRITE},
				},
				"storage": {
					"config": {authpkg.DELETE, authpkg.READ, authpkg.WRITE},
				},
			},
		},

		"permissions-merges-wildcard-resources-ditching-other-duplicates": {
			scopes: []authpkg.Scope{
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.READ, authpkg.WRITE}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.DELETE}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.WILDCARD}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.WILDCARD}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.DELETE}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.WILDCARD}},
			},
			expected: authpkg.Permissions{
				"identity": {
					"account":             {authpkg.WILDCARD},
					"account_credentials": {authpkg.WILDCARD},
				},
				"storage": {
					"config": {authpkg.WILDCARD},
				},
			},
		},

		"resource-wildcard": {
			scopes: []authpkg.Scope{
				{Service: "identity", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.WILDCARD}},

				{Service: "storage", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.WRITE}},
			},
			expected: authpkg.Permissions{
				"identity": {
					"*":                   {authpkg.READ},
					"account":             {authpkg.WRITE},
					"account_credentials": {authpkg.WILDCARD},
				},
				"storage": {
					"*":      {authpkg.READ},
					"config": {authpkg.WRITE},
				},
			},
		},

		"service-wildcard": {
			scopes: []authpkg.Scope{
				{Service: "*", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},

				{Service: "identity", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.WILDCARD}},

				{Service: "storage", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.WRITE}},
			},
			expected: authpkg.Permissions{
				"*": {
					"*": {authpkg.READ},
				},
				"identity": {
					"*":                   {authpkg.READ},
					"account":             {authpkg.WRITE},
					"account_credentials": {authpkg.WILDCARD},
				},
				"storage": {
					"*":      {authpkg.READ},
					"config": {authpkg.WRITE},
				},
			},
		},

		"wildcard-multiple-merges": {
			scopes: []authpkg.Scope{
				{Service: "*", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},

				{Service: "identity", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "identity", Resource: "account", Actions: []authpkg.Action{authpkg.WRITE}},
				{Service: "identity", Resource: "account_credentials", Actions: []authpkg.Action{authpkg.WILDCARD}},
				{Service: "identity", Resource: "*", Actions: []authpkg.Action{authpkg.WRITE}},

				{Service: "storage", Resource: "*", Actions: []authpkg.Action{authpkg.READ}},
				{Service: "storage", Resource: "config", Actions: []authpkg.Action{authpkg.WRITE}},

				{Service: "*", Resource: "*", Actions: []authpkg.Action{authpkg.WRITE}},
			},
			expected: authpkg.Permissions{
				"*": {
					"*": {authpkg.READ, authpkg.WRITE},
				},
				"identity": {
					"*":                   {authpkg.READ, authpkg.WRITE},
					"account":             {authpkg.WRITE},
					"account_credentials": {authpkg.WILDCARD},
				},
				"storage": {
					"*":      {authpkg.READ},
					"config": {authpkg.WRITE},
				},
			},
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {

				permissions := authpkg.NewPermissions(tt.scopes)
				assert.Equal(t, tt.expected, permissions)
			},
		)
	}
}

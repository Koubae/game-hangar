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
				Action:   "read",
			},
			expectedErr: nil,
		},
		"valid-scope-write": {
			scope: "identity:account:write",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Action:   "write",
			},
			expectedErr: nil,
		},
		"valid-scope-delete": {
			scope: "identity:account:delete",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Action:   "delete",
			},
			expectedErr: nil,
		},
		"valid-scope-spaces-and-uppercase-normalized": {
			scope: "IdenTity : ACCOUNT : DeletE   ",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Action:   "delete",
			},
			expectedErr: nil,
		},

		"invalid-scope-action-is-wildcard": {
			scope: "identity:account:*",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "account",
				Action:   "*",
			},
			expectedErr: nil,
		},
		"invalid-scope-resource-is-wildcard": {
			scope: "identity:*:read",
			expected: &authpkg.Scope{
				Service:  "identity",
				Resource: "*",
				Action:   "read",
			},
			expectedErr: nil,
		},
		"invalid-scope-service-is-wildcard": {
			scope: "*:account:read",
			expected: &authpkg.Scope{
				Service:  "*",
				Resource: "account",
				Action:   "read",
			},
			expectedErr: nil,
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

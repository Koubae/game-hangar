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
		"empty-scope": {
			scope:       "",
			expected:    nil,
			expectedErr: nil,
		},
		"empty-scope-with-spaces": {
			scope:       "                     ",
			expected:    nil,
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

func TestPermissions_ParsePermissions(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		scope       string
		expected    authpkg.Permissions
		expectedErr error
	}{
		"empty-scope": {
			scope:       "",
			expected:    authpkg.Permissions{},
			expectedErr: nil,
		},
		"empty-scope-with-spaces": {
			scope:       "                     ",
			expected:    authpkg.Permissions{},
			expectedErr: nil,
		},
		"permissions-parsed-default-scope-format": {
			scope: "identity:account:read,write|identity:account_credentials:write|storage:config:read|storage:storage:read,write",
			expected: authpkg.Permissions{
				"identity": {
					"account":             {authpkg.READ, authpkg.WRITE},
					"account_credentials": {authpkg.WRITE},
				},
				"storage": {
					"config":  {authpkg.READ},
					"storage": {authpkg.READ, authpkg.WRITE},
				},
			},
			expectedErr: nil,
		},
		"permissions-parsed-duplicate-scopes": {
			scope: "identity:account:read|identity:account:write|identity:account_credentials:write|storage:config:read|storage:storage:read|storage:storage:write",
			expected: authpkg.Permissions{
				"identity": {
					"account":             {authpkg.READ, authpkg.WRITE},
					"account_credentials": {authpkg.WRITE},
				},
				"storage": {
					"config":  {authpkg.READ},
					"storage": {authpkg.READ, authpkg.WRITE},
				},
			},
			expectedErr: nil,
		},
		"wildcard-multiple-merges": {
			scope: "*:*:read|identity:*:read|identity:account:write|identity:account_credentials:*|identity:*:write|storage:*:read|storage:config:write|*:*:write",
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

				permissions, err := authpkg.ParsePermissions(tt.scope)

				if tt.expectedErr != nil {
					assert.Nil(t, permissions)
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.expectedErr)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.expected, permissions)
			},
		)
	}
}

func TestPermissions_IsActionGranted(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		permissions authpkg.Permissions
		service     string
		resource    string
		action      authpkg.Action
		expected    bool
	}{
		"nil permissions returns false": {
			permissions: nil,
			service:     "identity",
			resource:    "account",
			action:      authpkg.READ,
			expected:    false,
		},
		"empty permissions returns false": {
			permissions: authpkg.Permissions{},
			service:     "identity",
			resource:    "account",
			action:      authpkg.READ,
			expected:    false,
		},
		"exact match grants allowed action": {
			permissions: authpkg.Permissions{
				"identity": {
					"account": {authpkg.READ, authpkg.WRITE},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.READ,
			expected: true,
		},
		"exact match denies missing action": {
			permissions: authpkg.Permissions{
				"identity": {
					"account": {authpkg.READ},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.WRITE,
			expected: false,
		},
		"resource wildcard grants action": {
			permissions: authpkg.Permissions{
				"identity": {
					"*": {authpkg.WRITE},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.WRITE,
			expected: true,
		},
		"service wildcard grants action": {
			permissions: authpkg.Permissions{
				"*": {
					"account": {authpkg.READ},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.READ,
			expected: true,
		},
		"global wildcard grants action": {
			permissions: authpkg.Permissions{
				"*": {
					"*": {authpkg.READ, authpkg.WRITE},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.READ,
			expected: true,
		},
		"global wildcard does not grant action because missing": {
			permissions: authpkg.Permissions{
				"*": {
					"*": {authpkg.READ, authpkg.WRITE},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.DELETE,
			expected: false,
		},
		"action wildcard grants any action": {
			permissions: authpkg.Permissions{
				"identity": {
					"account": {authpkg.WILDCARD},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.DELETE,
			expected: true,
		},
		"no permission for different service": {
			permissions: authpkg.Permissions{
				"identity": {
					"account": {authpkg.READ},
				},
			},
			service:  "storage",
			resource: "config",
			action:   authpkg.READ,
			expected: false,
		},
		"no permission for different resource": {
			permissions: authpkg.Permissions{
				"identity": {
					"account": {authpkg.READ},
				},
			},
			service:  "identity",
			resource: "account_credentials",
			action:   authpkg.READ,
			expected: false,
		},
		"no permission for different action": {
			permissions: authpkg.Permissions{
				"identity": {
					"account": {authpkg.READ},
				},
			},
			service:  "identity",
			resource: "account",
			action:   authpkg.WRITE,
			expected: false,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				got := tt.permissions.IsActionGranted(tt.service, tt.resource, tt.action)
				assert.Equal(t, tt.expected, got)
			},
		)
	}
}

func TestPermissions_Differance(t *testing.T) {
	tests := []struct {
		name        string
		p           authpkg.Permissions
		other       authpkg.Permissions
		wantDiff    authpkg.Permissions
		wantMissing []string
	}{
		{
			name: "all permissions match",
			p: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read", "write"},
				},
			},
			other: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
				},
			},
			wantDiff: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
				},
			},
			wantMissing: []string{},
		},
		{
			name: "missing service",
			p: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
				},
			},
			other: authpkg.Permissions{
				"svc2": {
					"res9": []authpkg.Action{"write"},
				},
			},
			wantDiff:    authpkg.PermissionEmpty,
			wantMissing: []string{"svc2:res9:write"},
		},
		{
			name: "missing resource",
			p: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
				},
			},
			other: authpkg.Permissions{
				"svc1": {
					"res2": []authpkg.Action{"write"},
				},
			},
			wantDiff:    authpkg.PermissionEmpty,
			wantMissing: []string{"svc1:res2:write"},
		},
		{
			name: "missing action",
			p: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
				},
			},
			other: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"write"},
				},
			},
			wantDiff:    authpkg.PermissionEmpty,
			wantMissing: []string{"svc1:res1:write"},
		},
		{
			name: "partial match and partial missing",
			p: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read", "write"},
				},
			},
			other: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read", "delete"},
				},
			},
			wantDiff: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
				},
			},
			wantMissing: []string{"svc1:res1:delete"},
		},
		{
			name: "multiple services and resources",
			p: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read", "write"},
					"res2": []authpkg.Action{"list"},
				},
				"svc2": {
					"resA": []authpkg.Action{"execute"},
				},
			},
			other: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read", "delete"},
					"res2": []authpkg.Action{"list"},
				},
				"svc2": {
					"resA": []authpkg.Action{"execute"},
					"resB": []authpkg.Action{"run"},
				},
			},
			wantDiff: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
					"res2": []authpkg.Action{"list"},
				},
				"svc2": {
					"resA": []authpkg.Action{"execute"},
				},
			},
			wantMissing: []string{
				"svc1:res1:delete",
				"svc2:resB:run",
			},
		},
		{
			name: "nil input returns empty",
			p:    nil,
			other: authpkg.Permissions{
				"svc1": {
					"res1": []authpkg.Action{"read"},
				},
			},
			wantDiff:    authpkg.PermissionEmpty,
			wantMissing: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotDiff, gotMissing := tt.p.Differance(tt.other)

				assert.Equal(t, tt.wantDiff, gotDiff)
				assert.ElementsMatch(t, tt.wantMissing, gotMissing)
			},
		)
	}
}

func TestPermissions_Scope(t *testing.T) {
	tests := []struct {
		name string
		p    authpkg.Permissions
		want string
	}{
		{
			name: "nil permissions",
			p:    nil,
			want: "",
		},
		{
			name: "single permission",
			p: authpkg.Permissions{
				"identity": {
					"account": []authpkg.Action{"read"},
				},
			},
			want: "identity:account:read",
		},
		{
			name: "multiple actions same resource",
			p: authpkg.Permissions{
				"identity": {
					"account": []authpkg.Action{"write", "read"},
				},
			},
			want: "identity:account:read,write",
		},
		{
			name: "multiple services resources and actions",
			p: authpkg.Permissions{
				"identity": {
					"account":             []authpkg.Action{"write", "read"},
					"account_credentials": []authpkg.Action{"write"},
				},
				"storage": {
					"config":  []authpkg.Action{"read"},
					"storage": []authpkg.Action{"write", "read"},
				},
			},
			want: "identity:account:read,write|identity:account_credentials:write|storage:config:read|storage:storage:read,write",
		},
		{
			name: "multiple services resources and actions with wildcards",
			p: authpkg.Permissions{
				"*": {
					"*": []authpkg.Action{"read"},
				},
				"identity": {
					"*":                   []authpkg.Action{"read"},
					"account":             []authpkg.Action{"write", "read"},
					"account_credentials": []authpkg.Action{"write"},
				},
				"storage": {
					"config":  []authpkg.Action{"read"},
					"storage": []authpkg.Action{"write", "read"},
				},
			},
			want: "*:*:read|identity:*:read|identity:account:read,write|identity:account_credentials:write|storage:config:read|storage:storage:read,write",
		},
		{
			name: "skips empty actions",
			p: authpkg.Permissions{
				"identity": {
					"account": []authpkg.Action{},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := tt.p.Scope()
				assert.Equal(t, tt.want, got)
			},
		)
	}
}

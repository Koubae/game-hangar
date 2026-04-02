package testunit

import (
	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/identity/account"
	authModels "github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/pkg/testutil"
)

const (
	AppPrefix = "TESTING_"
	EnvFile   = ".env.testing"
)

var (
	ProviderUsername = &authModels.Provider{
		ID:          1,
		Source:      "global",
		Type:        string(authModels.Username),
		DisplayName: "Username",
		Category:    "managed",
		Created:     testutil.Now,
		Updated:     testutil.Now,
	}

	ProviderEmail = &authModels.Provider{
		ID:          2,
		Source:      "global",
		Type:        string(authModels.Email),
		DisplayName: "Email",
		Category:    "managed",
		Created:     testutil.Now,
		Updated:     testutil.Now,
	}

	ProviderUsernameID = int64(1)

	StrongPassword     = "StrongPassword123!"
	StrongPasswordHash = "$2a$10$fSi7cP.8W9AbkmVwjO5lvuA3gdEKU97YjewosAUMmLjn1PDLozzkm"

	CredIDTest01   = int64(9999)
	UsernameTest01 = "unit-test-user-01"

	AccountIDTest01    = uuid.New()
	AccountIDTest01Str = AccountIDTest01.String()
	AccountEmail       = "unit-test@test.com"
	AccountTest01      = &account.Account{
		ID:       AccountIDTest01Str,
		Username: UsernameTest01,
		Email:    &AccountEmail,
		Disabled: false,
		Created:  testutil.Now,
		Updated:  testutil.Now,
	}

	AccountCredentialTest01 = &authModels.AccountCredential{
		ID:         1,
		Credential: UsernameTest01,
		AccountID:  AccountIDTest01,
		ProviderID: 1,
		Secret:     StrongPasswordHash,
	}

	AdminUsernameTest01     = "unit-test-admin-01"
	AdminAccountIDTest01    = uuid.New()
	AdminAccountIDTest01Str = AdminAccountIDTest01.String()
	AdminAccountEmail       = "admin-test@test.com"
	AdminAccountTest01      = &account.Account{
		ID:       AdminAccountIDTest01Str,
		Username: AdminUsernameTest01,
		Email:    &AdminAccountEmail,
		Disabled: false,
		Created:  testutil.Now,
		Updated:  testutil.Now,
	}

	AdminAccountCredentialTest01 = &authModels.AccountCredential{
		ID:         2,
		Credential: AdminUsernameTest01,
		AccountID:  AdminAccountIDTest01,
		ProviderID: 1,
		Secret:     StrongPasswordHash,
	}

	AdminAccountPermissions = []*authModels.Permission{
		{
			ID:       1,
			Service:  "identity",
			Resource: "auth",
			Action:   "read",
			Created:  testutil.Now,
		},
		{
			ID:       2,
			Service:  "identity",
			Resource: "auth",
			Action:   "write",
			Created:  testutil.Now,
		},
		{
			ID:       3,
			Service:  "identity",
			Resource: "auth",
			Action:   "delete",
			Created:  testutil.Now,
		},
		{
			ID:       4,
			Service:  "identity",
			Resource: "account",
			Action:   "read",
			Created:  testutil.Now,
		},
		{
			ID:       5,
			Service:  "identity",
			Resource: "account",
			Action:   "write",
			Created:  testutil.Now,
		},
		{
			ID:       6,
			Service:  "identity",
			Resource: "account",
			Action:   "delete",
			Created:  testutil.Now,
		},
		{
			ID:       7,
			Service:  "storage",
			Resource: "config",
			Action:   "read",
			Created:  testutil.Now,
		},
		{
			ID:       8,
			Service:  "storage",
			Resource: "setting",
			Action:   "read",
			Created:  testutil.Now,
		},
		{
			ID:       9,
			Service:  "storage",
			Resource: "setting",
			Action:   "write",
			Created:  testutil.Now,
		},
		{
			ID:       10,
			Service:  "leaderboard",
			Resource: "leaderboard",
			Action:   "read",
			Created:  testutil.Now,
		},
		{
			ID:       11,
			Service:  "chat",
			Resource: "*",
			Action:   "*",
			Created:  testutil.Now,
		},
	}
)

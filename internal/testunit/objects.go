package testunit

import (
	"github.com/google/uuid"
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

	AccountIDTest01    = uuid.New()
	AccountIDTest01Str = AccountIDTest01.String()

	StrongPassword = "StrongPassword123!"

	CredIDTest01   = int64(9999)
	UsernameTest01 = "unit-test-user-01"
)

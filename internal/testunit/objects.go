package testunit

import (
	authModels "github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/testutil"
)

const AppPrefix = "TESTING_"

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
)

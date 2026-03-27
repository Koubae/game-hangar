package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/tests/integration"
	"github.com/koubae/game-hangar/tests/testobj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accountRepo "github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	accountSrv "github.com/koubae/game-hangar/internal/identity/app/modules/account/service"
	authSrv "github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"
)

func TestAccountAuthService_RegisterByUsername(t *testing.T) {
	ctx, container, tearDown := integration.SetupTestIntegrationIdentity(t)
	defer tearDown(integration.ResetDB)

	testID := uuid.NewString()[:8]
	tests := []struct {
		id          string
		source      string
		username    string
		errExpected error
	}{
		{
			id:          "on-err-provider-not-exists",
			source:      "provider-source-does-not-exists",
			username:    "account-integration-01" + "-" + testID,
			errExpected: authSrv.ErrGetProvider,
		},
		{
			id:          "on-err-provider-is-disabled",
			source:      testobj.ProviderSourceDisabled,
			username:    "account-integration-01" + "-" + testID,
			errExpected: authSrv.ErrProviderIsDisabled,
		},

		{
			id:          "on-err-credential-exists",
			source:      testobj.ProviderSourceDefault,
			username:    testobj.CredentialAccount01,
			errExpected: accountSrv.ErrRegistrationCredExists,
		},

		{
			id:          "on-err-credential-is-empty",
			source:      testobj.ProviderSourceDefault,
			username:    "",
			errExpected: accountRepo.ErrUsernameRequired,
		},
		{
			id:          "account-and-credential-are-created",
			source:      testobj.ProviderSourceDefault,
			username:    "account-integration-01" + "-" + testID,
			errExpected: nil,
		},
	}

	accountAuthSrv := container.AccountAuthService(nil)
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			accountID, credID, err := accountAuthSrv.RegisterByUsername(
				ctx,
				tt.source,
				tt.username,
				testobj.PassHash,
			)

			if tt.errExpected != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errExpected)
				assert.Nil(t, accountID)
				assert.Nil(t, credID)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, accountID)
				assert.NotNil(t, credID)

				var accountIDFromDB string
				if err := container.DB().SelectOne(
					ctx,
					"SELECT id::text FROM account WHERE id = $1",
					accountID).Scan(&accountIDFromDB); err != nil {
					require.NoError(t, err)
				}

				var credIDFromDB int64
				var accountIDFK string
				if err := container.DB().SelectOne(
					ctx,
					"SELECT id, account_id::text FROM account_credentials WHERE id = $1",
					credID,
				).Scan(&credIDFromDB, &accountIDFK); err != nil {
					require.NoError(t, err)
				}

				assert.Equal(t, *accountID, accountIDFromDB)
				assert.Equal(t, *credID, credIDFromDB)
				assert.Equal(t, accountIDFromDB, accountIDFK)

			}
		})
	}
}

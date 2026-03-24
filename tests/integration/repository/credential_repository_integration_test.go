package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	accountID1 string
	username1  string
	email1     string
}

type credExpected struct {
	AccountID  string
	ProviderID int64
	Credential string
}

func createTestAccount(
	t *testing.T,
	ctx context.Context,
	connector *postgres.ConnectorPostgres,
) *testData {
	id := uuid.NewString()
	username := "user-itegration-" + id[:8]
	email := fmt.Sprintf("%s@integration.test", username)

	var err error
	var accountID uuid.UUID
	err = connector.SelectOne(ctx, `
		INSERT INTO account (id, username, email)
			VALUES ($1, $2, $3)
		RETURNING id 
	;`, id, username, email).Scan(&accountID)
	require.NoError(t, err)
	require.NotEqual(t, 0, accountID)

	var credentialID int
	err = connector.SelectOne(ctx, `
		INSERT INTO account_credentials (credential, account_id, provider_id, secret)
			VALUES ($1, $2, $3, encode(digest('pass', 'sha256'), 'hex'))
		RETURNING id 
	`, username, accountID, 1).Scan(&credentialID)
	require.NoError(t, err)
	require.NotEqual(t, 0, credentialID)

	return &testData{
		accountID1: id,
		username1:  username,
		email1:     email,
	}
}

func TestCredentialRepository_GetCredentialByProvider(t *testing.T) {
	ctx, connector, tearDown := integration.DBWithCleanup(t)
	defer tearDown(integration.ResetDB)

	testData := createTestAccount(t, ctx, connector)
	tests := []struct {
		id          string
		providerID  int64
		username    string
		expected    *credExpected
		errReturned error
	}{
		{
			id:         "record-is-found",
			providerID: 1,
			username:   testData.username1,
			expected: &credExpected{
				AccountID:  testData.accountID1,
				ProviderID: 1,
				Credential: testData.username1,
			},
			errReturned: nil,
		},
		{
			id:          "record-is-not-found",
			providerID:  1,
			username:    "does-not-exists",
			expected:    nil,
			errReturned: database.ErrNotFound,
		},
		{
			id:          "record-is-not-found-wrong-provider",
			providerID:  2,
			username:    testData.username1,
			expected:    nil,
			errReturned: database.ErrNotFound,
		},
	}

	repo := repository.NewCredentialRepository()
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			model, err := repo.GetCredentialByProvider(
				ctx,
				connector,
				tt.providerID,
				tt.username,
			)

			if tt.id == "unexpected-error" {
				connector.Shutdown()
			}

			if tt.errReturned != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.errReturned, err)
			} else {
				assert.NoError(t, err)
			}

			var result *credExpected
			if tt.expected != nil {
				result = &credExpected{
					AccountID:  model.AccountID.String(),
					ProviderID: model.ProviderID,
					Credential: model.Credential,
				}
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

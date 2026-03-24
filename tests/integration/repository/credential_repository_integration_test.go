package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
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

	print(testData.accountID1)

	repo := repository.NewCredentialRepository()
	model, err := repo.GetCredentialByProvider(ctx, connector, 1, testData.username1)

	assert.NoError(t, err)
	assert.Equal(t, testData.username1, model.Credential)
}

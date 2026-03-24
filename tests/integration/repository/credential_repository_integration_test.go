package repository

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/tests/integration"
	"github.com/stretchr/testify/require"
)

func TestCredentialRepository_GetCredentialByProvider(t *testing.T) {
	ctx, connector, tearDown := integration.DBWithCleanup(t)
	defer tearDown(integration.ResetDB)

	id := uuid.NewString()
	username := "user-itegration-" + id[:8]
	email := fmt.Sprintf("%s@integration.test", username)

	var accountID uuid.UUID
	err := connector.SelectOne(ctx, `
		INSERT INTO account (id, username, email)
			VALUES ($1, $2, $3)
		RETURNING id 
	;`, id, username, email).Scan(&accountID)
	require.NoError(t, err)
	require.NotEqual(t, 0, accountID)
}

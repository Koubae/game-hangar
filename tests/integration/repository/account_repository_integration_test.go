package repository

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account"
	"github.com/koubae/game-hangar/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	ctx, connector, tearDown := integration.DBWithCleanup(t, true)
	defer tearDown(integration.ResetDB)

	testID := uuid.NewString()[:8]
	username1 := "account-integration-01" + "-" + testID
	email1 := fmt.Sprintf("%s@integration.test", username1)
	tests := []struct {
		id          string
		params      account.NewAccount
		errReturned error
	}{
		{
			id: "account-created",
			params: account.NewAccount{
				Username: username1,
				Email:    &email1,
			},
			errReturned: nil,
		},
		{
			id: "err-duplicate-account",
			params: account.NewAccount{
				Username: username1,
				Email:    &email1,
			},
			errReturned: errs.ResourceDuplicate,
		},
		{
			id: "account-created-email-can-be-null",
			params: account.NewAccount{
				Username: "account-integration-02" + "-" + testID,
				Email:    nil,
			},
			errReturned: nil,
		},
	}

	repo := account.NewAccountRepository()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				id, err := repo.CreateAccount(ctx, connector, tt.params)
				if tt.errReturned != nil {
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.errReturned)
					assert.Nil(t, id)
				} else {
					assert.NoError(t, err)
					assert.NotEqual(t, nil, id)
				}
			},
		)
	}
}

func TestAccountRepository_GetAccount(t *testing.T) {
	ctx, connector, tearDown := integration.DBWithCleanup(t, true)
	defer tearDown(integration.ResetDB)

	testID := uuid.NewString()[:8]
	username1 := "account-integration-01" + "-" + testID
	email1 := fmt.Sprintf("%s@integration.test", username1)

	repo := account.NewAccountRepository()

	id1, err := repo.CreateAccount(
		ctx,
		connector,
		account.NewAccount{Username: username1, Email: &email1},
	)
	if err != nil {
		require.NoError(t, err)
	}

	type modelExpected struct {
		ID       string
		Username string
		Email    *string
	}

	tests := []struct {
		id          string
		accountID   string
		expected    *modelExpected
		errReturned error
	}{
		{
			id:        "record-is-found",
			accountID: *id1,
			expected: &modelExpected{
				ID:       *id1,
				Username: username1,
				Email:    &email1,
			},
			errReturned: nil,
		},
		{
			id:          "record-is-not-found",
			accountID:   uuid.NewString(),
			expected:    nil,
			errReturned: errs.ResourceNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				model, err := repo.GetAccount(ctx, connector, tt.accountID)
				if tt.errReturned != nil {
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.errReturned)
					assert.Nil(t, model)
				} else {
					assert.NoError(t, err)
				}

				var result *modelExpected
				if tt.expected != nil {
					result = &modelExpected{
						ID:       model.ID,
						Username: model.Username,
						Email:    model.Email,
					}
				}

				assert.Equal(t, tt.expected, result)
			},
		)
	}
}

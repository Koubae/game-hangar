package repository

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/tests/integration"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	ctx, connector, tearDown := integration.DBWithCleanup(t)
	defer tearDown(integration.ResetDB)

	testID := uuid.NewString()[:8]
	username1 := "account-integration-01" + "-" + testID
	email1 := fmt.Sprintf("%s@integration.test", username1)
	tests := []struct {
		id          string
		params      repository.NewAccount
		errReturned error
	}{
		{
			id: "account-created",
			params: repository.NewAccount{
				Username: username1,
				Email:    &email1,
			},
			errReturned: nil,
		},
		{
			id: "err-duplicate-account",
			params: repository.NewAccount{
				Username: username1,
				Email:    &email1,
			},
			errReturned: &database.ErrDuplicate{},
		},
		{
			id: "account-created-email-can-be-null",
			params: repository.NewAccount{
				Username: "account-integration-02" + "-" + testID,
				Email:    nil,
			},
			errReturned: nil,
		},
	}

	repo := repository.NewAccountRepository()
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			id, err := repo.CreateAccount(ctx, connector, tt.params)
			if tt.errReturned != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.errReturned, err)
				assert.Nil(t, id)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, nil, id)
			}
		})
	}
}

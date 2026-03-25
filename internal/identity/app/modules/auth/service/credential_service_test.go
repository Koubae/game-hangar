package service

import (
	"context"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredentialService_GetCredentialByProvider(t *testing.T) {
	t.Parallel()

	testunit.Setup()
	connector := testunit.MockDBConnector()
	providerID := int64(1)
	username := "unit-test-user-123"
	tests := []struct {
		id            string
		provider      int64
		credential    string
		setupMock     func(repo *testunit.MockCredentialRepository)
		expected      *string
		errorReturned error
	}{
		{
			id:         "record-is-found",
			provider:   providerID,
			credential: username,
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"GetCredentialByProvider",
						mock.Anything,
						connector,
						providerID,
						username,
					).
					Return(&model.AccountCredential{
						ID:         1,
						Credential: username,
						AccountID:  testutil.AccountIDTest01,
						ProviderID: 1,
					}, nil).
					Once()
			},
			expected:      &username,
			errorReturned: nil,
		},
		{
			id:         "record-is-not-found",
			provider:   providerID,
			credential: username,
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"GetCredentialByProvider",
						mock.Anything,
						connector,
						providerID,
						username,
					).
					Return(nil, database.ErrNotFound).
					Once()
			},
			expected:      nil,
			errorReturned: database.ErrNotFound,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			repo := new(testunit.MockCredentialRepository)
			tt.setupMock(repo)

			service := NewCredentialService(connector, repo)

			result, err := service.GetCredentialByProvider(
				ctx,
				tt.provider,
				tt.credential,
			)

			if tt.errorReturned != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errorReturned)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, *tt.expected, result.Credential)
			}
		})
	}
}

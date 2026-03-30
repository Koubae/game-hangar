package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredentialService_GetCredentialByProvider(t *testing.T) {
	t.Parallel()

	container := testunit.NewTestIdentityAppContainer(t)
	connector := container.DB()

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
					Return(
						&model.AccountCredential{
							ID:         1,
							Credential: username,
							AccountID:  testutil.AccountIDTest01,
							ProviderID: 1,
						}, nil,
					).
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
					Return(nil, errs.ResourceNotFound).
					Once()
			},
			expected:      nil,
			errorReturned: errs.ResourceNotFound,
		},
		{
			id:         "on-db-error",
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
					Return(nil, testunit.ErrDBGeneric).
					Once()
			},
			expected:      nil,
			errorReturned: testunit.ErrDBGeneric,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				repo := container.CredentialRepository().(*testunit.MockCredentialRepository)
				tt.setupMock(repo)

				_service := service.NewCredentialService(connector, repo)

				result, err := _service.GetCredentialByProvider(
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
			},
		)
	}
}

func TestCredentialService_CreateCredentialTypeUsername(t *testing.T) {
	t.Parallel()

	testunit.Setup()
	connector := testunit.MockDBConnector()
	username := "unit-test-user-123"

	tests := []struct {
		id            string
		credential    string
		accountID     uuid.UUID
		provider      *model.Provider
		setupMock     func(repo *testunit.MockCredentialRepository)
		expected      int64
		errorReturned error
	}{
		{
			id:         "credential-created",
			credential: username,
			accountID:  testutil.AccountIDTest01,
			provider:   testunit.ProviderUsername,
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"CreateAccountCredential",
						mock.Anything,
						connector,
						mock.AnythingOfType("repository.NewAccountCredential"),
					).
					Run(
						func(args mock.Arguments) {
							params := args.Get(2).(repository.NewAccountCredential)
							assert.Equal(t, "password", params.SecretType)
							assert.True(
								t,
								params.Verified,
								"CreateCredentialTypeUsername should set verified to true",
							)
						},
					).
					Return(int64(9999), nil).
					Once()
			},
			expected:      int64(9999),
			errorReturned: nil,
		},
		{
			id:         "credential-is-duplicated",
			credential: username,
			accountID:  testutil.AccountIDTest01,
			provider:   testunit.ProviderUsername,
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"CreateAccountCredential",
						mock.Anything,
						connector,
						mock.AnythingOfType("repository.NewAccountCredential"),
					).
					Run(
						func(args mock.Arguments) {
							params := args.Get(2).(repository.NewAccountCredential)
							assert.Equal(t, "password", params.SecretType)
							assert.True(
								t,
								params.Verified,
								"CreateCredentialTypeUsername should set verified to true",
							)
						},
					).
					Return(int64(0), &database.ErrDuplicate{}).
					Once()
			},
			expected:      int64(0),
			errorReturned: &database.ErrDuplicate{},
		},
		{
			id:         "on-db-error",
			credential: username,
			accountID:  testutil.AccountIDTest01,
			provider:   testunit.ProviderUsername,
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"CreateAccountCredential",
						mock.Anything,
						connector,
						mock.AnythingOfType("repository.NewAccountCredential"),
					).
					Run(
						func(args mock.Arguments) {
							params := args.Get(2).(repository.NewAccountCredential)
							assert.Equal(t, "password", params.SecretType)
							assert.True(
								t,
								params.Verified,
								"CreateCredentialTypeUsername should set verified to true",
							)
						},
					).
					Return(int64(0), testunit.ErrDBGeneric).
					Once()
			},
			expected:      int64(0),
			errorReturned: testunit.ErrDBGeneric,
		},

		{
			id:         "on-wrong-provider-type",
			credential: username,
			accountID:  testutil.AccountIDTest01,
			provider:   testunit.ProviderEmail,
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.AssertNotCalled(t, "CreateAccountCredential")
			},
			expected:      int64(0),
			errorReturned: service.ErrCreateCredentialIncorrectProviderType,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				t.Parallel()

				repo := new(testunit.MockCredentialRepository)
				tt.setupMock(repo)

				service := service.NewCredentialService(connector, repo)

				result, err := service.CreateCredentialTypeUsername(
					ctx,
					tt.credential,
					tt.accountID,
					tt.provider,
					"secret-hash-sha256",
				)

				if tt.errorReturned != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tt.errorReturned)
					assert.Equal(t, int64(0), result)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, result)
				}
			},
		)
	}
}

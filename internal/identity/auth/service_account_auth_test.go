package auth_test

import (
	"context"
	"testing"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"github.com/koubae/game-hangar/tests/testobj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountAuthService_RegisterByUsernameProviderErr(t *testing.T) {
	t.Parallel()

	type expected struct {
		accountID *string
		credID    *int64
	}

	tests := []struct {
		id          string
		source      string
		credential  string
		setupMock   func(repo *testunit.MockProviderRepository)
		expected    *expected
		errExpected error
	}{
		{
			id:         "err-on-provider-not-exists",
			source:     "provider-does-not-exists",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockProviderRepository) {
				repo.
					On(
						"GetProvider",
						mock.Anything,
						mock.Anything,
						"provider-does-not-exists",
						"username",
					).
					Return(nil, errspkg.ResourceNotFound).
					Once()
			},
			expected:    nil,
			errExpected: errs.ProviderNotFound,
		},
		{
			id:         "err-on-provider-disabled",
			source:     "global",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockProviderRepository) {
				repo.
					On(
						"GetProvider",
						mock.Anything,
						mock.Anything,
						"global",
						"username",
					).
					Return(&auth.Provider{ID: 1, Source: "global", Type: "username", Disabled: true}, nil).
					Once()
			},
			expected:    nil,
			errExpected: errs.ProviderDisabled,
		},
	}

	container := testunit.NewTestIdentityAppContainer(t)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				t.Parallel()

				repo := container.ProviderRepository().(*testunit.MockProviderRepository)
				tt.setupMock(repo)

				service := container.AccountAuthService(nil)
				accountID, credID, err := service.RegisterByUsername(
					ctx,
					tt.source,
					tt.credential,
					testobj.PassHash,
				)

				if tt.errExpected != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tt.errExpected)
					assert.Nil(t, accountID)
					assert.Nil(t, credID)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, *tt.expected, expected{accountID: accountID, credID: credID})
				}
			},
		)
	}
}

func TestAccountAuthService_RegisterByUsernameCredentialErr(t *testing.T) {
	t.Parallel()

	type expected struct {
		accountID *string
		credID    *int64
	}

	tests := []struct {
		id          string
		source      string
		credential  string
		setupMock   func(repo *testunit.MockCredentialRepository)
		expected    *expected
		errExpected error
	}{
		{
			id:         "on-err-credential-exists",
			source:     "global",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"GetCredentialByProvider",
						mock.Anything,
						mock.Anything,
						int64(1),
						"test-cred",
					).
					Return(
						&auth.AccountCredential{
							ID:         1,
							Credential: "test-cred",
							AccountID:  testunit.AccountIDTest01,
							ProviderID: 1,
						}, nil,
					).
					Once()
			},
			expected:    nil,
			errExpected: errs.AccountCredDuplicate,
		},
		{
			id:         "on-err-credential-generic-db-err",
			source:     "global",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"GetCredentialByProvider",
						mock.Anything,
						mock.Anything,
						int64(1),
						"test-cred",
					).
					Return(nil, errspkg.DBError).
					Once()
			},
			expected:    nil,
			errExpected: errspkg.DBError,
		},
	}

	container := testunit.NewTestIdentityAppContainer(t)
	providerRepo := container.ProviderRepository().(*testunit.MockProviderRepository)
	providerRepo.On(
		"GetProvider",
		mock.Anything,
		mock.Anything,
		"global",
		"username",
	).
		Return(&auth.Provider{ID: 1, Source: "global", Type: "username", Disabled: false}, nil)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				t.Parallel()

				repo := container.CredentialRepository().(*testunit.MockCredentialRepository)
				tt.setupMock(repo)

				service := container.AccountAuthService(nil)
				accountID, credID, err := service.RegisterByUsername(
					ctx,
					tt.source,
					tt.credential,
					testobj.PassHash,
				)

				if tt.errExpected != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tt.errExpected)
					assert.Nil(t, accountID)
					assert.Nil(t, credID)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, *tt.expected, expected{accountID: accountID, credID: credID})
				}
			},
		)
	}
}

func TestAccountAuthService_RegisterByUsernameAccountAndCredentialCreation(
	t *testing.T,
) {
	t.Parallel()

	type expected struct {
		accountID *string
		credID    *int64
	}

	accountID := testunit.AccountIDTest01.String()
	credIDExpected := int64(9999)
	tests := []struct {
		id            string
		source        string
		credential    string
		setupMock     func(repo *testunit.MockAccountRepository)
		setupMockCred func(repo *testunit.MockCredentialRepository)
		expected      *expected
		errExpected   error
	}{
		{
			id:         "on-err-credential-exists",
			source:     "global",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockAccountRepository) {
				repo.
					On(
						"CreateAccount",
						mock.Anything,
						mock.Anything,
						account.NewAccount{
							Username: "test-cred",
							Email:    nil,
						},
					).
					Return(nil, errspkg.DBError).
					Once()
			},
			setupMockCred: nil,
			expected:      nil,
			errExpected:   errspkg.DBError,
		},
		{
			id:         "account-credential-created",
			source:     "global",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockAccountRepository) {
				repo.
					On(
						"CreateAccount",
						mock.Anything,
						mock.Anything,
						account.NewAccount{
							Username: "test-cred",
							Email:    nil,
						},
					).
					Return(&accountID, nil).
					Once()
			},
			setupMockCred: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"CreateAccountCredential",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).
					Return(credIDExpected, nil).
					Once()
			},
			expected: &expected{
				accountID: &accountID,
				credID:    &credIDExpected,
			},
			errExpected: nil,
		},

		{
			id:         "on-credential-creation-err",
			source:     "global",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockAccountRepository) {
				repo.
					On(
						"CreateAccount",
						mock.Anything,
						mock.Anything,
						account.NewAccount{
							Username: "test-cred",
							Email:    nil,
						},
					).
					Return(&accountID, nil).
					Once()
			},
			setupMockCred: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"CreateAccountCredential",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).
					Return(nil, errspkg.DBError).
					Once()
			},
			expected:    nil,
			errExpected: errspkg.DBError,
		},
		{
			id:         "account-credential-created",
			source:     "global",
			credential: "test-cred",
			setupMock: func(repo *testunit.MockAccountRepository) {
				repo.
					On(
						"CreateAccount",
						mock.Anything,
						mock.Anything,
						account.NewAccount{
							Username: "test-cred",
							Email:    nil,
						},
					).
					Return(&accountID, nil).
					Once()
			},
			setupMockCred: func(repo *testunit.MockCredentialRepository) {
				repo.
					On(
						"CreateAccountCredential",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).
					Return(credIDExpected, nil).
					Once()
			},
			expected: &expected{
				accountID: &accountID,
				credID:    &credIDExpected,
			},
			errExpected: nil,
		},
	}

	container := testunit.NewTestIdentityAppContainer(t)
	providerRepo := container.ProviderRepository().(*testunit.MockProviderRepository)
	providerRepo.On(
		"GetProvider",
		mock.Anything,
		mock.Anything,
		"global",
		"username",
	).
		Return(&auth.Provider{ID: 1, Source: "global", Type: "username", Disabled: false}, nil)

	credentialRepo := container.CredentialRepository().(*testunit.MockCredentialRepository)
	credentialRepo.On(
		"GetCredentialByProvider",
		mock.Anything,
		mock.Anything,
		int64(1),
		"test-cred",
	).
		Return(nil, errspkg.ResourceNotFound)

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				repo := container.AccountRepository().(*testunit.MockAccountRepository)
				tt.setupMock(repo)

				if tt.setupMockCred != nil {
					tt.setupMockCred(credentialRepo)
				}

				service := container.AccountAuthService(nil)
				accountID, credID, err := service.RegisterByUsername(
					ctx,
					tt.source,
					tt.credential,
					testobj.PassHash,
				)

				if tt.errExpected != nil {
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.errExpected)
					assert.Nil(t, accountID)
					assert.Nil(t, credID)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, *tt.expected, expected{accountID: accountID, credID: credID})
				}
			},
		)
	}
}

package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredentialRepository_GetCredentialByProvider(t *testing.T) {
	t.Parallel()

	providerID := int64(1)
	username := "unit-test-user-123"
	tests := []struct {
		id          string
		provider    int64
		username    string
		expected    *model.AccountCredential
		errThrown   error
		errReturned error
	}{
		{
			id:       "record-is-found",
			provider: providerID,
			username: username,
			expected: &model.AccountCredential{
				ID:         1,
				Credential: username,
				AccountID:  testutil.AccountIDTest01,
				ProviderID: 1,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
				Disabled:   false,
				DisabledAt: nil,
			},
			errThrown:   nil,
			errReturned: nil,
		},
		{
			id:          "record-is-not-found",
			provider:    providerID,
			username:    username,
			expected:    nil,
			errThrown:   pgx.ErrNoRows,
			errReturned: database.ErrNotFound,
		},
	}

	modelToValues := func(s *model.AccountCredential) []any {
		if s == nil {
			return []any{}
		}
		return []any{
			s.ID,
			s.Credential,
			s.AccountID,
			s.ProviderID,
			s.Secret,
			s.SecretType,
			s.Verified,
			s.VerifiedAt,
			s.Disabled,
			s.DisabledAt,
		}
	}

	fieldsCount := reflect.TypeFor[model.AccountCredential]().NumField()
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			common.CreateLogger(common.LogLevelError, "")
			mockRow := new(testutil.MockRow)
			mockRow.MockScan(
				fieldsCount,
				tt.errThrown,
				modelToValues(tt.expected)...,
			)

			mockPool := new(testutil.MockDBPool)
			mockPool.On("QueryRow", mock.Anything, mock.Anything, providerID, username).Return(mockRow)

			connector := postgres.ConnectorPostgres{Pool: mockPool}
			repo := NewCredentialRepository()

			model, err := repo.GetCredentialByProvider(
				context.Background(),
				&connector,
				providerID,
				username,
			)

			if tt.errThrown != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errReturned)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, model)
			mockPool.AssertExpectations(t)
		})
	}
}

func TestCredentialRepository_CreateAccountCredential(t *testing.T) {
	t.Parallel()

	params := NewAccountCredential{
		Credential: "unit-test-user-123",
		AccountID:  testutil.AccountIDTest01,
		ProviderID: 1,
		Secret:     "sha255-secret",
		SecretType: "password",
		Verified:   true,
		VerifiedAt: &testutil.Now,
	}
	expectedID := int64(1234)

	common.CreateLogger(common.LogLevelError, "")
	mockRow := new(testutil.MockRow)
	mockRow.MockScan(1, nil, expectedID)

	mockPool := new(testutil.MockDBPool)
	mockPool.On("QueryRow", mock.Anything, mock.Anything, pgx.StrictNamedArgs{
		"credential":  params.Credential,
		"account_id":  params.AccountID,
		"provider_id": params.ProviderID,
		"secret":      params.Secret,
		"secret_type": params.SecretType,
		"verified":    params.Verified,
		"verified_at": params.VerifiedAt,
	}).Return(mockRow)

	ctx := context.Background()
	connector := postgres.ConnectorPostgres{Pool: mockPool}
	repo := NewCredentialRepository()

	id, err := repo.CreateAccountCredential(ctx, &connector, params)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)
	assert.Equal(t, expectedID, id)
}

func TestCredentialRepository_CreateAccountCredentialOnErrors(t *testing.T) {
	t.Parallel()

	providerID := int64(1)
	username := "unit-test-user-123"
	mockedDBErr := errors.New("mocked-db-error")
	tests := []struct {
		id          string
		params      *NewAccountCredential
		expectedID  int64
		errThrown   error
		errReturned error
	}{
		{
			id: "validation-err-verified-required",
			params: &NewAccountCredential{
				Credential: username,
				AccountID:  testutil.AccountIDTest01,
				ProviderID: providerID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: nil,
			},
			expectedID:  int64(0),
			errThrown:   nil,
			errReturned: ErrVerifiedAtRequired,
		},
		{
			id: "validation-err-nil-when-f",
			params: &NewAccountCredential{
				Credential: username,
				AccountID:  testutil.AccountIDTest01,
				ProviderID: providerID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   false,
				VerifiedAt: &testutil.Now,
			},
			expectedID:  int64(0),
			errThrown:   nil,
			errReturned: ErrVerifiedNilWhenIsFalse,
		},
		{
			id: "on-db-error-any",
			params: &NewAccountCredential{
				Credential: username,
				AccountID:  testutil.AccountIDTest01,
				ProviderID: providerID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedID:  int64(0),
			errThrown:   mockedDBErr,
			errReturned: mockedDBErr,
		},
		{
			id: "on-db-error-duplicate-resource",
			params: &NewAccountCredential{
				Credential: username,
				AccountID:  testutil.AccountIDTest01,
				ProviderID: providerID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedID:  int64(0),
			errThrown:   testutil.DBMockErrDuplicateKey,
			errReturned: database.ErrrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			common.CreateLogger(common.LogLevelError, "")
			mockRow := new(testutil.MockRow)
			mockRow.MockScan(1, tt.errThrown, tt.expectedID)

			params := tt.params

			mockPool := new(testutil.MockDBPool)
			mockPool.On("QueryRow", mock.Anything, mock.Anything, pgx.StrictNamedArgs{
				"credential":  params.Credential,
				"account_id":  params.AccountID,
				"provider_id": params.ProviderID,
				"secret":      params.Secret,
				"secret_type": params.SecretType,
				"verified":    params.Verified,
				"verified_at": params.VerifiedAt,
			}).Return(mockRow)

			ctx := context.Background()
			connector := postgres.ConnectorPostgres{Pool: mockPool}
			repo := NewCredentialRepository()

			id, err := repo.CreateAccountCredential(ctx, &connector, *params)

			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.errReturned)

			assert.Equal(t, tt.expectedID, id)
		})
	}
}

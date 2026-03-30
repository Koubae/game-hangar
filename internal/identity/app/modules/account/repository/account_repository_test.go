package repository_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	t.Parallel()

	emailTest := "test@unit.test"
	mockedDBErr := errors.New("mocked-db-error")
	tests := []struct {
		id          string
		params      *repository.NewAccount
		expected    string
		errThrown   error
		errReturned error
	}{
		{
			id: "resource-created",
			params: &repository.NewAccount{
				Username: "account-01",
				Email:    &emailTest,
			},
			expected:    "account-01",
			errThrown:   nil,
			errReturned: nil,
		},
		{
			id: "on-db-error-duplicate-resource",
			params: &repository.NewAccount{
				Username: "account-01",
				Email:    &emailTest,
			},
			expected:    "",
			errThrown:   testutil.DBMockErrDuplicateKey,
			errReturned: errs.ResourceDuplicate,
		},
		{
			id: "on-db-error-any",
			params: &repository.NewAccount{
				Username: "account-01",
				Email:    &emailTest,
			},
			expected:    "",
			errThrown:   mockedDBErr,
			errReturned: mockedDBErr,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				common.CreateLogger(common.LogLevelDPanic, "")
				mockRow := new(testutil.MockRow)
				mockRow.MockScan(1, tt.errThrown, tt.expected)

				params := tt.params

				mockPool := new(testutil.MockDBPool)
				mockPool.On(
					"QueryRow", mock.Anything, mock.Anything, pgx.StrictNamedArgs{
						"username": params.Username,
						"email":    params.Email,
					},
				).
					Return(mockRow)

				connector := postgres.ConnectorPostgres{Pool: mockPool}
				repo := repository.NewAccountRepository()

				id, err := repo.CreateAccount(ctx, &connector, *params)

				if tt.errReturned != nil {
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.errReturned)
					assert.Nil(t, id)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, &tt.expected, id)
				}
			},
		)
	}
}

func TestAccountRepository_GetAccount(t *testing.T) {
	t.Parallel()

	emailTest := "test@unit.test"
	tests := []struct {
		id          string
		accountID   string
		expected    *account.Account
		errThrown   error
		errReturned error
	}{
		{
			id:        "record-is-found",
			accountID: "06e1b677-a4fe-42cf-8afd-ceec867d1fa5",
			expected: &account.Account{
				ID:       "06e1b677-a4fe-42cf-8afd-ceec867d1fa5",
				Username: "account-01",
				Email:    &emailTest,
				Disabled: false,
				Created:  testutil.Now,
				Updated:  testutil.Now,
			},
			errThrown:   nil,
			errReturned: nil,
		},
		{
			id:          "record-is-not-found",
			accountID:   "06e1b677-a4fe-42cf-8afd-ceec867d1fa5",
			expected:    nil,
			errThrown:   pgx.ErrNoRows,
			errReturned: errs.ResourceNotFound,
		},
	}

	modelToValues := func(s *account.Account) []any {
		if s == nil {
			return []any{}
		}
		return []any{
			s.ID,
			s.Username,
			s.Email,
			s.Disabled,
			s.Created,
			s.Updated,
		}
	}

	ctx := context.Background()
	fieldsCount := reflect.TypeFor[account.Account]().NumField()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				common.CreateLogger(common.LogLevelDPanic, "")
				mockRow := new(testutil.MockRow)
				mockRow.MockScan(
					fieldsCount,
					tt.errThrown,
					modelToValues(tt.expected)...,
				)

				mockPool := new(testutil.MockDBPool)
				mockPool.On("QueryRow", mock.Anything, mock.Anything, pgx.StrictNamedArgs{"id": tt.accountID}).
					Return(mockRow)

				connector := postgres.ConnectorPostgres{Pool: mockPool}
				repo := repository.NewAccountRepository()

				_model, err := repo.GetAccount(ctx, &connector, tt.accountID)

				if tt.errThrown != nil {
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.errReturned)
				} else {
					assert.NoError(t, err)
				}

				assert.Equal(t, tt.expected, _model)
				mockPool.AssertExpectations(t)
			},
		)
	}
}

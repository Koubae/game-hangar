package testunit

import (
	"context"
	"errors"

	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/mock"
)

var ErrDBGeneric = errors.New("mock-db-error")

func MockDBConnector() *postgres.ConnectorPostgres {
	mockPool := new(testutil.MockDBPool)

	mockPool.On("BeginTx", mock.Anything, mock.Anything).
		Return(testutil.DefaultStubPgxTx, nil)
	connector := postgres.ConnectorPostgres{Pool: mockPool}
	return &connector
}

type MockProviderRepository struct {
	mock.Mock
}

func NewMockProviderRepository() auth.IProviderRepository {
	return new(MockProviderRepository)
}

func (m *MockProviderRepository) LoadProviders(
	ctx context.Context,
	db database.DBTX,
) {
	_ = m.Called(ctx, db)
}

func (m *MockProviderRepository) GetProvider(
	ctx context.Context,
	db database.DBTX,
	source string,
	_type string,
) (*auth.Provider, error) {
	args := m.Called(ctx, db, source, _type)

	provider, _ := args.Get(0).(*auth.Provider)
	return provider, args.Error(1)
}

type MockCredentialRepository struct {
	mock.Mock
}

func NewMockCredentialRepository() auth.ICredentialRepository {
	return new(MockCredentialRepository)
}

func (m *MockCredentialRepository) GetCredentialByProvider(
	ctx context.Context,
	db database.DBTX,
	providerID int64,
	credential string,
) (*auth.AccountCredential, error) {
	args := m.Called(ctx, db, providerID, credential)

	model, _ := args.Get(0).(*auth.AccountCredential)
	return model, args.Error(1)
}

func (m *MockCredentialRepository) CreateAccountCredential(
	ctx context.Context,
	db database.DBTX,
	params auth.NewAccountCredential,
) (int64, error) {
	args := m.Called(ctx, db, params)

	id, _ := args.Get(0).(int64)
	return id, args.Error(1)
}

type MockAccountRepository struct {
	mock.Mock
}

func NewMockAccountRepository() account.IAccountRepository {
	return new(MockAccountRepository)
}

func (m *MockAccountRepository) CreateAccount(
	ctx context.Context,
	db database.DBTX,
	params account.NewAccount,
) (*string, error) {
	args := m.Called(ctx, db, params)

	id, _ := args.Get(0).(*string)
	return id, args.Error(1)
}

func (m *MockAccountRepository) GetAccount(
	ctx context.Context,
	db database.DBTX,
	id string,
) (*account.Account, error) {
	args := m.Called(ctx, db, id)
	return args.Get(0).(*account.Account), args.Error(1)
}

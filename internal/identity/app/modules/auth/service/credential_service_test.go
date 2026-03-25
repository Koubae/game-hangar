package service

import (
	"context"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/stretchr/testify/mock"
)

type MockCredentialRepository struct {
	mock.Mock
}

func (m *MockCredentialRepository) GetCredentialByProvider(
	ctx context.Context,
	db database.DBTX,
	providerID int64,
	credential string,
) (*model.AccountCredential, error) {
	args := m.Called(ctx, db, providerID, credential)

	model, _ := args.Get(0).(*model.AccountCredential)
	return model, args.Error(1)
}

func (m *MockCredentialRepository) CreateAccountCredential(
	ctx context.Context,
	db database.DBTX,
	params repository.NewAccountCredential,
) (int64, error) {
	args := m.Called(ctx, db, params)

	id, _ := args.Get(0).(int64)
	return id, args.Error(1)
}

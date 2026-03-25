package testunit

import (
	"context"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/stretchr/testify/mock"
)

type MockProviderRepository struct {
	mock.Mock
}

func (m *MockProviderRepository) LoadProviders(ctx context.Context, db database.DBTX) {
	_ = m.Called(ctx, db)
}

func (m *MockProviderRepository) GetProvider(ctx context.Context, db database.DBTX, source string, _type string) (*model.Provider, error) {
	args := m.Called(ctx, db, source, _type)

	provider, _ := args.Get(0).(*model.Provider)
	return provider, args.Error(1)
}

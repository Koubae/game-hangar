package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPool) Close() {
	m.Called()
}

func newTestProviderRepository(mockPool *MockPool, cache map[string]*model.Provider) *ProviderRepository {
	return &ProviderRepository{
		DBConnector: &postgres.ConnectorPostgres{
			Pool: mockPool,
		},
		providersCache: cache,
	}
}

func TestProviderRepository_GetProvider_CacheHit(t *testing.T) {
	t.Parallel()

	now := time.Now()

	expected := &model.Provider{
		ID:          1,
		Name:        "stream",
		DisplayName: "Steam",
		Category:    "platform",
		Disabled:    false,
		Created:     now,
		Updated:     now,
	}

	mockPool := new(MockPool)
	repo := newTestProviderRepository(mockPool, map[string]*model.Provider{
		"steam": expected,
	})

	got, err := repo.GetProvider(context.Background(), "steam")

	assert.NoError(t, err)
	assert.Same(t, expected, got)
	mockPool.AssertExpectations(t)

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func TestProviderRepository_GetProvider_CacheMiss(t *testing.T) {
	t.Parallel()

	mockPool := new(MockPool)
	repo := newTestProviderRepository(mockPool, map[string]*model.Provider{
		"email": {
			ID:          2,
			Name:        "email",
			DisplayName: "Email",
			Category:    "internal",
			Disabled:    false,
		},
	})

	got, err := repo.GetProvider(context.Background(), "steam")

	assert.NoError(t, err)
	assert.Nil(t, got)
	mockPool.AssertExpectations(t)
}

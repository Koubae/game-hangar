package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	mockPool := new(testutil.MockDBPool)
	connector := postgres.ConnectorPostgres{Pool: mockPool}

	repo := &ProviderRepository{providersCache: map[string]*model.Provider{"steam": expected}}
	got, err := repo.GetProvider(context.Background(), &connector, "steam")

	assert.NoError(t, err)
	assert.Same(t, expected, got)
	mockPool.AssertExpectations(t)

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func TestProviderRepository_GetProvider_CacheMiss(t *testing.T) {
	t.Parallel()

	common.CreateLogger(common.LogLevelInfo, "")
	expected := &model.Provider{
		ID:          2,
		Name:        "email",
		DisplayName: "Email",
		Category:    "internal",
		Disabled:    false,
	}

	mockRow := new(testutil.MockRow)
	mockRow.MockScan(7, nil, expected.ID, expected.Name, expected.DisplayName, expected.Category, expected.Disabled)

	mockPool := new(testutil.MockDBPool)
	mockPool.On("QueryRow", mock.Anything, mock.Anything, "steam").Return(mockRow)

	connector := postgres.ConnectorPostgres{Pool: mockPool}
	repo := &ProviderRepository{providersCache: map[string]*model.Provider{"email": expected}}
	got, err := repo.GetProvider(context.Background(), &connector, "steam")

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	mockPool.AssertExpectations(t)
}

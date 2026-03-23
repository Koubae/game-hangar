package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
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

	tests := []struct {
		id       string
		expected *model.Provider
		err      error
	}{
		{
			id: "record-is-found",
			expected: &model.Provider{
				ID:          2,
				Name:        "email",
				DisplayName: "Email",
				Category:    "internal",
				Disabled:    false,
			},
			err: nil,
		},

		{
			id:       "record-is-not-found",
			expected: nil,
			err:      database.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			common.CreateLogger(common.LogLevelInfo, "")
			mockRow := new(testutil.MockRow)

			if tt.expected != nil {
				mockRow.MockScan(7, nil,
					tt.expected.ID,
					tt.expected.Name,
					tt.expected.DisplayName,
					tt.expected.Category,
					tt.expected.Disabled,
				)
			} else {
				mockRow.MockScan(7, pgx.ErrNoRows)
			}

			mockPool := new(testutil.MockDBPool)
			mockPool.On("QueryRow", mock.Anything, mock.Anything, "steam").Return(mockRow)

			connector := postgres.ConnectorPostgres{Pool: mockPool}
			repo := &ProviderRepository{providersCache: map[string]*model.Provider{"email": tt.expected}}
			got, err := repo.GetProvider(context.Background(), &connector, "steam")

			if tt.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, got)
			mockPool.AssertExpectations(t)
		})
	}
}

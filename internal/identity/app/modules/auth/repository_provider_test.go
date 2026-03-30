package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProviderRepository_GetProvider_CacheHit(t *testing.T) {
	t.Parallel()

	now := time.Now()

	expected := &auth.Provider{
		ID:          1,
		Source:      "global",
		Type:        "stream",
		DisplayName: "Steam",
		Category:    "platform",
		Disabled:    false,
		Created:     now,
		Updated:     now,
	}

	mockPool := new(testutil.MockDBPool)
	connector := postgres.ConnectorPostgres{Pool: mockPool}

	repo := &auth.ProviderRepository{
		ProvidersCache: map[string]map[string]*auth.Provider{
			"global": {"steam": expected},
		},
	}
	got, err := repo.GetProvider(
		context.Background(),
		&connector,
		"global",
		"steam",
	)

	assert.NoError(t, err)
	assert.Same(t, expected, got)
	mockPool.AssertExpectations(t)
}

func TestProviderRepository_GetProvider_CacheMiss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id       string
		expected *auth.Provider
		err      error
	}{
		{
			id: "record-is-found",
			expected: &auth.Provider{
				ID:          2,
				Source:      "global",
				Type:        "email",
				DisplayName: "Email",
				Category:    "internal",
				Disabled:    false,
			},
			err: nil,
		},

		{
			id:       "record-is-not-found",
			expected: nil,
			err:      errs.ResourceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				common.CreateLogger(common.LogLevelDPanic, "")
				mockRow := new(testutil.MockRow)

				if tt.expected != nil {
					mockRow.MockScan(
						8, nil,
						tt.expected.ID,
						tt.expected.Source,
						tt.expected.Type,
						tt.expected.DisplayName,
						tt.expected.Category,
						tt.expected.Disabled,
					)
				} else {
					mockRow.MockScan(8, pgx.ErrNoRows)
				}

				mockPool := new(testutil.MockDBPool)
				mockPool.On(
					"QueryRow", mock.Anything, mock.Anything, pgx.StrictNamedArgs{
						"source": "global", "type": "steam",
					},
				).
					Return(mockRow)

				connector := postgres.ConnectorPostgres{Pool: mockPool}

				repo := &auth.ProviderRepository{
					ProvidersCache: map[string]map[string]*auth.Provider{
						"global": {"email": tt.expected},
					},
				}
				got, err := repo.GetProvider(
					context.Background(),
					&connector,
					"global",
					"steam",
				)

				if tt.err != nil {
					assert.Error(t, err)
					assert.ErrorIs(t, err, tt.err)
				} else {
					assert.NoError(t, err)
				}

				assert.Equal(t, tt.expected, got)
				mockPool.AssertExpectations(t)
			},
		)
	}
}

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProviderRepository struct {
	mock.Mock
}

func (m *MockProviderRepository) LoadProviders(ctx context.Context, db database.DBTX) {
	_ = m.Called(ctx, db)
}

func (m *MockProviderRepository) GetProvider(ctx context.Context, db database.DBTX, name string) (*model.Provider, error) {
	args := m.Called(ctx, db, name)

	provider, _ := args.Get(0).(*model.Provider)
	return provider, args.Error(1)
}

func TestProviderService_IsProviderEnabled(t *testing.T) {
	t.Parallel()

	common.CreateLogger("INFO", "/tmp/")

	tests := []struct {
		name         string
		providerName string
		setupMock    func(repo *MockProviderRepository)
		want         bool
	}{
		{
			name:         "returns true when provider is enabled",
			providerName: "steam",
			setupMock: func(repo *MockProviderRepository) {
				repo.
					On("GetProvider", mock.Anything, mock.Anything, "steam").
					Run(func(args mock.Arguments) {
						ctx := args.Get(0).(context.Context)
						_, hasDeadline := ctx.Deadline()
						assert.True(t, hasDeadline, "expected context to have deadline")
					}).
					Return(&model.Provider{
						Name:     "steam",
						Disabled: false,
					}, nil).
					Once()
			},
			want: true,
		},
		{
			name:         "returns false when provider is disabled",
			providerName: "steam",
			setupMock: func(repo *MockProviderRepository) {
				repo.
					On("GetProvider", mock.Anything, mock.Anything, "steam").
					Run(func(args mock.Arguments) {
						ctx := args.Get(0).(context.Context)
						_, hasDeadline := ctx.Deadline()
						assert.True(t, hasDeadline, "expected context to have deadline")
					}).
					Return(&model.Provider{
						Name:     "steam",
						Disabled: true,
					}, nil).
					Once()
			},
			want: false,
		},
		{
			name:         "returns false when repository returns error",
			providerName: "steam",
			setupMock: func(repo *MockProviderRepository) {
				repo.
					On("GetProvider", mock.Anything, mock.Anything, "steam").
					Run(func(args mock.Arguments) {
						ctx := args.Get(0).(context.Context)
						_, hasDeadline := ctx.Deadline()
						assert.True(t, hasDeadline, "expected context to have deadline")
					}).
					Return((*model.Provider)(nil), errors.New("repository failure")).
					Once()
			},
			want: false,
		},
	}

	mockPool := new(testutil.MockDBPool)
	connector := postgres.ConnectorPostgres{Pool: mockPool}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(MockProviderRepository)
			tt.setupMock(repo)

			svc := NewProviderService(&connector, repo)

			got := svc.IsProviderEnabled(context.Background(), tt.providerName)

			assert.Equal(t, tt.want, got)
			repo.AssertExpectations(t)
		})
	}
}

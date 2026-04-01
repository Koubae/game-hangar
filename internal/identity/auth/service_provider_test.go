package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProviderService_IsProviderEnabled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		id     string
		source string
		_type  string

		setupMock func(repo *testunit.MockProviderRepository)
		want      bool
	}{
		{
			id:     "returns true when provider is enabled",
			source: "global",
			_type:  "steam",
			setupMock: func(repo *testunit.MockProviderRepository) {
				repo.
					On(
						"GetProvider",
						mock.Anything,
						mock.Anything,
						"global",
						"steam",
					).
					Run(
						func(args mock.Arguments) {
							ctx := args.Get(0).(context.Context)
							_, hasDeadline := ctx.Deadline()
							assert.True(
								t,
								hasDeadline,
								"expected context to have deadline",
							)
						},
					).
					Return(
						&auth.Provider{
							Source:   "global",
							Type:     "steam",
							Disabled: false,
						}, nil,
					).
					Once()
			},
			want: true,
		},
		{
			id:     "returns false when provider is disabled",
			source: "global",
			_type:  "steam",
			setupMock: func(repo *testunit.MockProviderRepository) {
				repo.
					On(
						"GetProvider",
						mock.Anything,
						mock.Anything,
						"global",
						"steam",
					).
					Run(
						func(args mock.Arguments) {
							ctx := args.Get(0).(context.Context)
							_, hasDeadline := ctx.Deadline()
							assert.True(
								t,
								hasDeadline,
								"expected context to have deadline",
							)
						},
					).
					Return(
						&auth.Provider{
							Source:   "global",
							Type:     "steam",
							Disabled: true,
						}, nil,
					).
					Once()
			},
			want: false,
		},
		{
			id:     "returns false when repository returns error",
			source: "global",
			_type:  "steam",
			setupMock: func(repo *testunit.MockProviderRepository) {
				repo.
					On("GetProvider", mock.Anything, mock.Anything, "global", "steam").
					Run(
						func(args mock.Arguments) {
							ctx := args.Get(0).(context.Context)
							_, hasDeadline := ctx.Deadline()
							assert.True(
								t,
								hasDeadline,
								"expected context to have deadline",
							)
						},
					).
					Return((*auth.Provider)(nil), errors.New("repository failure")).
					Once()
			},
			want: false,
		},
	}

	container := testunit.NewTestIdentityAppContainer(t)
	connector := container.DB()
	for _, tt := range tests {
		t.Run(
			tt.id, func(t *testing.T) {
				repo := container.ProviderRepository().(*testunit.MockProviderRepository)
				tt.setupMock(repo)

				svc := auth.NewProviderService(connector, repo)

				got := svc.IsProviderEnabled(
					context.Background(),
					tt.source,
					tt._type,
				)

				assert.Equal(t, tt.want, got)
				repo.AssertExpectations(t)
			},
		)
	}
}

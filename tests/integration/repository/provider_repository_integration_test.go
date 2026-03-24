package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/tests/integration"
)

func setupTest(t *testing.T) *postgres.ConnectorPostgres {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_ = common.NewConfig(common.CreateLogger(common.LogLevelInfo, ""), ".env.testing", integration.AppPrefix)

	config, err := postgres.LoadConfig(integration.AppPrefix)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	connector, err := postgres.NewConnector(config)
	if err != nil {
		t.Fatalf("Failed to create connector (is database running?): %v", err)
	}

	return connector
}

func TestProviderRepository_GetProvider(t *testing.T) {
	connector := setupTest(t)
	defer connector.Shutdown()

	tests := []struct {
		id       string
		name     string
		expected model.Provider
	}{
		{
			id:   "provider-username",
			name: "username",
			expected: model.Provider{
				Name:        "username",
				DisplayName: "Username",
				Category:    "managed",
			},
		},
		{
			id:   "provider-email",
			name: "email",
			expected: model.Provider{
				Name:        "email",
				DisplayName: "Email",
				Category:    "managed",
			},
		},
		{
			id:   "provider-device",
			name: "device",
			expected: model.Provider{
				Name:        "device",
				DisplayName: "Device",
				Category:    "managed",
			},
		},
		{
			id:   "provider-guest",
			name: "guest",
			expected: model.Provider{
				Name:        "guest",
				DisplayName: "Guest",
				Category:    "anonymous",
			},
		},
		{
			id:   "provider-anonymous",
			name: "anonymous",
			expected: model.Provider{
				Name:        "anonymous",
				DisplayName: "Anonymous",
				Category:    "anonymous",
			},
		},
		{
			id:   "provider-steam",
			name: "steam",
			expected: model.Provider{
				Name:        "steam",
				DisplayName: "Steam",
				Category:    "platform",
			},
		},
		{
			id:   "provider-playstation",
			name: "psn",
			expected: model.Provider{
				Name:        "psn",
				DisplayName: "PlayStation Network",
				Category:    "platform",
			},
		},
		{
			id:   "provider-xbox",
			name: "xbox",
			expected: model.Provider{
				Name:        "xbox",
				DisplayName: "Xbox",
				Category:    "platform",
			},
		},
		{
			id:   "provider-nintendo",
			name: "nintendo",
			expected: model.Provider{
				Name:        "nintendo",
				DisplayName: "Nintendo",
				Category:    "platform",
			},
		},
	}

	providerRepository := repository.NewProviderRepository()
	providerRepository.LoadProviders(context.Background(), connector)
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			provider, err := providerRepository.GetProvider(context.Background(), connector, tt.name)
			if err != nil {
				t.Fatalf("Failed to get provider: %v", err)
			}
			if provider == nil {
				t.Fatalf("Provider is nil")
			}

			if provider.Name != tt.expected.Name {
				t.Fatalf("Provider name is not '%s' got: %s\n ", tt.expected.Name, provider.Name)
			}
			if provider.DisplayName != tt.expected.DisplayName {
				t.Fatalf("Provider display name is not '%s' \n got: %s", tt.expected.DisplayName, provider.DisplayName)
			}
			if provider.Category != tt.expected.Category {
				t.Fatalf("Provider category is not '%s' \n got: %s", tt.expected.Category, provider.Category)
			}
		})
	}
}

func TestProviderRepository_GetProviderFoundWhenCacheMiss(t *testing.T) {
	connector := setupTest(t)
	defer connector.Shutdown()

	providerName := "username"
	providerRepository := repository.NewProviderRepository()
	provider, err := providerRepository.GetProvider(context.Background(), connector, providerName)
	if err != nil {
		t.Fatalf("Failed to get provider: %v", err)
	}
	if provider == nil {
		t.Fatalf("Provider is nil")
	}

	if provider.Name != providerName {
		t.Fatalf("Provider name is not '%s' got: %s\n ", providerName, provider.Name)
	}
	if provider.DisplayName != "Username" {
		t.Fatalf("Provider display name is not '%s' \n got: %s", "Username", provider.DisplayName)
	}
	if provider.Category != "managed" {
		t.Fatalf("Provider category is not '%s' \n got: %s", "managed", provider.Category)
	}
}

func TestProviderRepository_GetProviderNotFound(t *testing.T) {
	connector := setupTest(t)
	defer connector.Shutdown()

	providerNotExists := "not-exists"
	providerRepository := repository.NewProviderRepository()
	providerRepository.LoadProviders(context.Background(), connector)

	provider, err := providerRepository.GetProvider(context.Background(), connector, providerNotExists)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			t.Fatalf("Failed to get provider: %v", err)
		}
	}
	if provider != nil {
		t.Fatalf("Provider is not nil")
	}
}

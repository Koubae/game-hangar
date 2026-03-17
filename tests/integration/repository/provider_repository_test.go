package repository

import (
	"context"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
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

	providerRepository := repository.NewProviderRepository(connector)
	provider, err := providerRepository.GetProvider(context.Background(), "username")
	if err != nil {
		t.Fatalf("Failed to get provider: %v", err)
	}
	if provider == nil {
		t.Fatalf("Provider is nil")
	}

	if provider.Name != "username" {
		t.Fatalf("Provider name is not %s got: %s\n ", "username", provider.Name)
	}
	if provider.DisplayName != "Username" {
		t.Fatalf("Provider display name is not %s \n got: %s", "Username", provider.DisplayName)
	}
	if provider.Category != "managed" {
		t.Fatalf("Provider category is not %s \n got: %s", "managed", provider.Category)
	}
	if provider.Disabled != false {
		t.Fatalf("Provider disabled is not %t \n got: %t", false, provider.Disabled)
	}
}

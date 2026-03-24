package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
)

const AppPrefix = "TESTING_"

func SetupTest(t *testing.T) *postgres.ConnectorPostgres {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_ = common.NewConfig(common.CreateLogger(common.LogLevelInfo, ""), ".env.testing", AppPrefix)

	config, err := postgres.LoadConfig(AppPrefix)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	connector, err := postgres.NewConnector(config)
	if err != nil {
		t.Fatalf("Failed to create connector (is database running?): %v", err)
	}

	return connector
}

func ResetDB(ctx context.Context, connector *postgres.ConnectorPostgres) error {
	tables := strings.Join([]string{
		`"public"."account_credentials"`,
		`"public"."account"`,
	}, ", ")

	query := `TRUNCATE TABLE ` + tables + ` RESTART IDENTITY CASCADE`
	if _, err := connector.SQL(ctx, query); err != nil {
		return fmt.Errorf("truncate test tables: %w", err)
	}

	return nil
}

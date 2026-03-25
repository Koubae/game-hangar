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

type (
	DBTearDownFN    func(ctx context.Context, connector *postgres.ConnectorPostgres) error
	DBTearDownTasks func(tasks ...DBTearDownFN)
)

func DBWithCleanup(
	t *testing.T,
) (context.Context, *postgres.ConnectorPostgres, DBTearDownTasks) {
	t.Helper()

	ctx := context.Background()
	connector := IntegrationTestConnector(t)

	return ctx, connector, func(tasks ...DBTearDownFN) {
		t.Helper()

		t.Cleanup(func() {
			defer connector.Shutdown()

			for _, fn := range tasks {
				if fn == nil {
					continue
				}

				err := fn(ctx, connector)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func IntegrationTestConnector(t *testing.T) *postgres.ConnectorPostgres {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_ = common.NewConfig(
		common.CreateLogger(common.LogLevelDPanic, ""),
		".env.testing",
		AppPrefix,
	)

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

package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	identityContainer "github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/tests/testobj"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const AppPrefix = "TESTING_"

type (
	DBTearDownFN    func(ctx context.Context, connector *postgres.ConnectorPostgres) error
	DBTearDownTasks func(tasks ...DBTearDownFN)
)

func SetupTestIntegrationIdentity(
	t *testing.T,
) (context.Context, *identityContainer.AppContainer, DBTearDownTasks) {
	MarkTestSlow(t)

	_, logger := LoadNewConfig()

	ctx, connector, tearDown := DBWithCleanup(t, false)

	dependencies, err := identityContainer.LoadAppDependenciesWithDefaultFactories(
		logger,
		connector,
	)
	require.NoError(t, err)

	container, err := identityContainer.NewAppContainer(
		AppPrefix,
		logger,
		dependencies,
	)
	require.NoError(t, err)

	return ctx, container, func(tasks ...DBTearDownFN) {
		t.Helper()

		t.Cleanup(
			func() {
				container.Logger().
					Info("Server has shutdown, cleaning up resources ...")

				if container != nil {
					if err := container.Shutdown(); err != nil {
						container.Logger().
							Error("Container Shutdown Failed", zap.Error(err))
					}
				}

				container.Logger().
					Info("Resource cleanup completed, terminating process...")
			},
		)

		tearDown(tasks...)
	}
}

func DBWithCleanup(
	t *testing.T, loadConfig bool,
) (context.Context, *postgres.ConnectorPostgres, DBTearDownTasks) {
	MarkTestSlow(t)

	if loadConfig {
		_, _ = LoadNewConfig()
	}

	ctx := context.Background()
	connector := integrationTestConnector(t)

	return ctx, connector, func(tasks ...DBTearDownFN) {
		t.Helper()

		t.Cleanup(
			func() {
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
			},
		)
	}
}

func integrationTestConnector(t *testing.T) *postgres.ConnectorPostgres {
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
	tables := strings.Join(
		[]string{
			`"public"."account_credentials"`,
			`"public"."account"`,
		}, ", ",
	)

	query := `TRUNCATE TABLE ` + tables + ` RESTART IDENTITY CASCADE`
	if _, err := connector.SQL(ctx, query); err != nil {
		return fmt.Errorf("truncate test tables: %w", err)
	}

	// NOTE: Re-Create Demo Data TODO: Let's improve this???
	if _, err := connector.SQL(ctx, testobj.SQLAccountDemoData); err != nil {
		return fmt.Errorf("error re-create account demo data, error: %w", err)
	}
	if _, err := connector.SQL(ctx, testobj.SQLCredentialsDemoData); err != nil {
		return fmt.Errorf(
			"error re-create credentials demo data, error: %w",
			err,
		)
	}

	return nil
}

func LoadNewConfig() (*common.Config, common.Logger) {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	config := common.NewConfig(loggerTmp, ".env.testing", AppPrefix)
	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	return config, logger
}

func MarkTestSlow(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
}

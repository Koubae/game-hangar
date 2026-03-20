package integration

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
)

func setupTest(t *testing.T) *postgres.ConnectorPostgres {
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

func TestIntegration_Connection(t *testing.T) {
	connector := setupTest(t)
	if connector == nil {
		t.Fatal("Connector is nil")
	}
	if connector.Pool == nil {
		t.Fatal("Pool is nil")
	}
}

func TestIntegration_Ping(t *testing.T) {
	connector := setupTest(t)
	err := connector.Ping(context.Background())
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}
}

func TestIntegration_SelectQuery(t *testing.T) {
	connector := setupTest(t)

	// In integration tests, we can cast to *pgxpool.Pool to access its methods
	// The Pool field in ConnectorPostgres is a PoolInterface

	type queryable interface {
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	}

	pool, ok := connector.Pool.(queryable)
	if !ok {
		// Let's try to see what it actually is
		t.Fatalf("Pool does not support QueryRow. Type: %T", connector.Pool)
	}

	var result int
	err := pool.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		t.Errorf("SELECT 1 query failed: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

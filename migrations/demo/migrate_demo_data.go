package main

import (
	"embed"
	"flag"

	"github.com/koubae/game-hangar/migrations"
	"go.uber.org/zap"
)

const (
	MigrationTable = "demo_data_migrations"
)

//go:embed sql/*.sql
var sqlMigrations embed.FS

var (
	limit     int
	action    string = "status"
	envFile   string = ".env"
	appPrefix string = "IDENTITY_"
)

/*
Usage:

	go run ./migrations/identity/migrate_demo_data.go -action status
	migrate_demo_data.go -action status
	migrate_demo_data.go -action up -limit <limit>
	migrate_demo_data.go -action down -limit <limit>

	-- Tests -- integration tests
	go run ./migrations/demo/migrate_demo_data.go -action status -env .env.testing -appPrefix TESTING_
	go run ./migrations/demo/migrate_demo_data.go -action up -limit 0 -env .env.testing -appPrefix TESTING_
	go run ./migrations/demo/migrate_demo_data.go -action down -limit 0 -env .env.testing -appPrefix TESTING_
*/
func main() {
	flag.StringVar(&action, "action", "status", "status|up|down")
	flag.IntVar(&limit, "limit", 0, "max number of migrations to apply (up/down only)")
	flag.StringVar(&envFile, "env", ".env", "environment file to use")
	flag.StringVar(&appPrefix, "appPrefix", "IDENTITY_", "application prefix")

	flag.Parse()

	migrator := migrations.InitializeMigrations(envFile, appPrefix, MigrationTable, sqlMigrations, false)
	defer migrator.Close()

	migrator.Logger.Info("Running migration (demo data): ", zap.String("action", action), zap.String("envFile", envFile), zap.String("appPrefix", appPrefix))

	result, err := migrator.Run(action, limit)
	if err != nil {
		migrator.Logger.Fatal("failed to run migration: ", zap.Error(err), zap.String("result", result))
	}
	migrator.Logger.Info("Migration result (demo data): ", zap.String("result", result))
}

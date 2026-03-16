package main

import (
	"embed"
	"flag"

	"github.com/koubae/game-hangar/migrations"
	"go.uber.org/zap"
)

const (
	AppPrefix = "IDENTITY_"
)

//go:embed sql/*.sql
var sqlMigrations embed.FS

var (
	action string = "status"
	limit  int
)

/*
Usage:

	go run ./migrations/identity/migrate_identity.go -action status
	migrate-identity -action status
	migrate-identity -action up <limit>
	migrate-identity -action down <limit>
*/
func main() {
	flag.StringVar(&action, "action", "status", "status|up|down")
	flag.IntVar(&limit, "limit", 0, "max number of migrations to apply (up/down only)")

	flag.Parse()

	migrator := migrations.InitializeMigrations(AppPrefix, sqlMigrations)
	defer migrator.Close()
	migrator.Logger.Info("Running migration: ", zap.String("action", action))

	result, err := migrator.Run(action, limit)
	if err != nil {
		migrator.Logger.Fatal("failed to run migration: ", zap.Error(err), zap.String("result", result))
	}
	migrator.Logger.Info("Migration result: ", zap.String("result", result))

}

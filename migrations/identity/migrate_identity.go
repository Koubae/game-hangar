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
var action string

/*
Usage:

	migrate-identity -action status
	migrate-identity -action up
	migrate-identity -action down
*/
func main() {
	flag.StringVar(&action, "action", "status", "status|up|down")
	flag.Parse()

	migrator := migrations.InitializeMigrations(AppPrefix, sqlMigrations)

	migrator.Logger.Info("Running migration: ", zap.String("action", action))

	result, err := migrator.Run(action)
	if err != nil {
		migrator.Logger.Fatal("failed to run migration: ", zap.Error(err))
	}
	migrator.Logger.Info("Migration result: ", zap.String("result", result))

}

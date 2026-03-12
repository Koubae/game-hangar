package run

import (
	"context"
	"log"

	app2 "github.com/koubae/game-hangar/internal/app"
)

func RunServer() {
	app := app2.NewApp()
	app.Start(context.Background())
	if err := app.Stop(); err != nil {
		log.Fatalf("Error while shutting down the server, error: %s", err)
	}
}

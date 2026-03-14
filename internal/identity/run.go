package identity

import (
	"context"
	"log"

	"github.com/koubae/game-hangar/internal/identity/app"
)

const AppPrefix = "IDENTITY_"

func RunServer() {
	application := app.NewApp(AppPrefix)
	application.Start(context.Background())
	if err := application.Stop(); err != nil {
		log.Fatalf("Error while shutting down the server, error: %s", err)
	}
}

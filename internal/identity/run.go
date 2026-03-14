package identity

import (
	"context"
	"log"

	"github.com/koubae/game-hangar/internal/identity/app"
	"github.com/koubae/game-hangar/pkg/web"
)

const AppPrefix = "IDENTITY_"

func RunServer() {
	application := web.NewHTTPApp(AppPrefix, web.Router, app.RouterRegister)
	application.Start(context.Background())
	if err := application.Stop(); err != nil {
		log.Fatalf("Error while shutting down the server, error: %s", err)
	}
}

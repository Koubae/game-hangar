package identity

import (
	"context"
	"log"

	"github.com/koubae/game-hangar/internal/identity/app"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/di"
	"github.com/koubae/game-hangar/pkg/web"
)

const AppPrefix = "IDENTITY_"

func RunServer() {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	config := common.NewConfig(loggerTmp, ".env", AppPrefix)
	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	container, err := di.NewContainer(AppPrefix, logger)
	if err != nil {
		log.Fatalf("Error while creating container, error: %s", err)
	}

	application := web.NewHTTPApp(AppPrefix, container, config, web.Router, app.RouterRegister(container))
	application.Start(context.Background())
	if err := application.Stop(); err != nil {
		log.Fatalf("Error while shutting down the server, error: %s", err)
	}
}

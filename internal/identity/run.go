package identity

import (
	"context"
	"log"

	"github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/web"
)

const AppPrefix = "IDENTITY_"

func RunServer() {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	config := common.NewConfig(loggerTmp, ".env", AppPrefix)
	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	_container, err := container.NewAppContainer(AppPrefix, logger, nil)
	if err != nil {
		log.Fatalf("Error while creating container, error: %s", err)
	}

	providerRepository := _container.ProviderRepository()
	providerRepository.LoadProviders(context.Background(), _container.DB())

	application := web.NewHTTPApp(
		AppPrefix,
		_container,
		config,
		web.Router,
		RouterRegister(_container),
	)
	application.Start(context.Background())
	if err := application.Stop(); err != nil {
		log.Fatalf("Error while shutting down the server, error: %s", err)
	}
}

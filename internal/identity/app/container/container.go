package container

import (
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityContainer interface {
	di.Container

	DBConnector() *postgres.ConnectorPostgres

	// Repositories
	ProviderRepository() repository.IProviderRepository
}

type AppContainer struct {
	logger    common.Logger
	connector *postgres.ConnectorPostgres

	// Repositories
	providerRepository repository.IProviderRepository
}

func NewAppContainer(appPrefix string, logger common.Logger) (*AppContainer, error) {
	dbConfig, err := postgres.LoadConfig(appPrefix)
	if err != nil {
		return nil, err
	}
	connector, err := postgres.InitConnector(dbConfig)
	if err != nil {
		return nil, err
	}

	logger.Info("database connection established", zap.String("db", connector.String()))

	providerRepository := repository.NewProviderRepository(connector)

	return &AppContainer{
		logger:             logger,
		connector:          connector,
		providerRepository: providerRepository,
	}, nil
}

// ------------------------------------------
//
//	Implements di.Container interface
//
// ------------------------------------------
func (c *AppContainer) Logger() common.Logger {
	return c.logger
}

func (c *AppContainer) Shutdown() error {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	defer func() {
		if z, ok := c.Logger().(*common.AppLogger); ok {
			z.LogCloser(loggerTmp, z.Logger)
		}
	}()

	c.DBConnector().Shutdown()
	c.Logger().Info("database connection closed")
	return nil
}

// ------------------------------------------
// 	Implements IdentityContainer interface
// ------------------------------------------

func (c *AppContainer) DBConnector() *postgres.ConnectorPostgres {
	return c.connector
}

func (c *AppContainer) ProviderRepository() *repository.IProviderRepository {
	return &c.providerRepository

}

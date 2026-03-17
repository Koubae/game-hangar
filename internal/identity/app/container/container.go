package container

import (
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityContainer interface {
	di.Container

	DB() *postgres.ConnectorPostgres
}

type AppContainer struct {
	logger common.Logger
	db     *postgres.ConnectorPostgres

	// Repositories
}

func NewAppContainer(appPrefix string, logger common.Logger) (*AppContainer, error) {
	dbConfig, err := postgres.LoadConfig(appPrefix)
	if err != nil {
		return nil, err
	}
	db, err := postgres.NewConnector(dbConfig)
	if err != nil {
		return nil, err
	}

	logger.Info("database connection established", zap.String("db", db.String()))

	return &AppContainer{
		logger: logger,
		db:     db,
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

	c.DB().Shutdown()
	c.Logger().Info("database connection closed")
	return nil
}

// ------------------------------------------
// 	Implements IdentityContainer interface
// ------------------------------------------

func (c *AppContainer) DB() *postgres.ConnectorPostgres {
	return c.db
}

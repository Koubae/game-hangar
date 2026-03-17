package container

import (
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"go.uber.org/zap"
)

type AppContainer struct {
	logger common.Logger
	DB     *postgres.ConnectorPostgres
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
		DB:     db,
	}, nil
}

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

	if c.DB != nil {
		c.DB.Shutdown()
		c.Logger().Info("database connection closed")
	}
	return nil
}

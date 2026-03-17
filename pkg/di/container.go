package di

import (
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"go.uber.org/zap"
)

type Container struct {
	DB     *postgres.ConnectorPostgres
	Logger common.Logger
}

func NewContainer(appPrefix string, logger common.Logger) (*Container, error) {
	dbConfig, err := postgres.LoadConfig(appPrefix)
	if err != nil {
		return nil, err
	}
	db, err := postgres.NewConnector(dbConfig)
	if err != nil {
		return nil, err
	}

	logger.Info("database connection established", zap.String("db", db.String()))

	return &Container{
		DB:     db,
		Logger: logger,
	}, nil
}

func (c *Container) Shutdown() error {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	defer func() {
		if z, ok := c.Logger.(*common.AppLogger); ok {
			z.LogCloser(loggerTmp, z.Logger)
		}
	}()

	if c.DB != nil {
		c.DB.Shutdown()
		c.Logger.Info("database connection closed")
	}
	return nil
}

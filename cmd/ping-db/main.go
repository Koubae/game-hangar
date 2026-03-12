package main

import (
	"github.com/koubae/game-hangar/internal/app/settings"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"go.uber.org/zap"
)

// Tests Database connection
func main() {
	config := settings.NewConfig(common.CreateLogger(common.LogLevelInfo, ""))
	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	logger.Info("ping-db script initialized... ")

	dbConfig, err := postgres.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load database configuration", zap.Error(err))
	}
	dbPool, err := postgres.NewConnector(dbConfig)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	} else if dbPool == nil {
		logger.Fatal("Database connection pool is nil...")
	}
	defer dbPool.Shutdown()

	logger.Info("database connection established... ", zap.String("dbConfig", dbConfig.String()))

}

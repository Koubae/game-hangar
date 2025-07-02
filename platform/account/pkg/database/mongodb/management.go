package mongodb

import (
	"errors"
	"github.com/koubae/game-hangar/account/pkg/utils"
)

type DatabaseConfig struct {
	Driver string
	Uri    string
	DBName string
}

var databaseConfig *DatabaseConfig

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	databaseDriver := utils.GetEnvString("APP_DATABASE_1_DRIVER", "")
	databaseUri := utils.GetEnvString("APP_DATABASE_1_URI", "")
	name := utils.GetEnvString("APP_DATABASE_1_NAME", "")
	if databaseDriver == "" || databaseUri == "" || name == "" {
		return nil, errors.New("database configuration missing")
	}

	config := &DatabaseConfig{
		Driver: databaseDriver,
		Uri:    databaseUri,
		DBName: name,
	}

	databaseConfig = config
	return config, nil
}

func GetDatabaseConfig() *DatabaseConfig {
	if databaseConfig == nil {
		panic("DatabaseConfig is not initialized, call LoadDatabaseConfig() first!")
	}
	return databaseConfig
}

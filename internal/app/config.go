package app

import (
	"fmt"
	"slices"

	"github.com/joho/godotenv"
	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

type Environment string

const (
	EnvTesting Environment = "testing"
	EnvDev     Environment = "dev"
	EnvStating Environment = "staging"
	EnvProd    Environment = "prod"
)

var Environments = [4]Environment{EnvTesting, EnvDev, EnvStating, EnvProd}

type Config struct {
	Host        string
	Port        int
	AppName     string
	AppVersion  string
	Env         Environment
	LogLevel    common.LogLevel
	LogFilePath string
}

func (c Config) GetAppURL() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		panic("Config is not initialized. Please initialize the Config before using it.")
	}
	return config
}

func NewConfig(logger *zap.Logger) *Config {
	_ = godotenv.Load(".env")

	host := common.GetEnvString("APP_HOST", "")
	port := common.GetEnvInt("APP_PORT", 8080)

	env := Environment(common.GetEnvString("APP_ENV", string(EnvDev)))
	if !slices.Contains(Environments[:], env) {
		logger.Panic(fmt.Sprintf("Invalid APP_ENV: %s, supported envs are: %v", env, Environments))
	}

	logLevel := common.LogLevel(common.GetEnvString("APP_LOG_LEVEL", string(common.LogLevelInfo)))
	if !slices.Contains(common.LogLevels[:], logLevel) {
		logger.Panic(fmt.Sprintf("Invalid LOG_LEVEL: %s, supported levels are: %v", logLevel, common.LogLevels))
	}
	logFilePath := common.GetEnvString("APP_LOG_FILE", "logs/app.log")

	config = &Config{
		Host:        host,
		Port:        port,
		AppName:     "game-hangar",
		Env:         env,
		LogLevel:    logLevel,
		LogFilePath: logFilePath,
	}
	return config
}

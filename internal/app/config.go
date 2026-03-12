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
	AppName    string
	AppVersion string
	Env        Environment
	// server configs
	Host                       string
	Port                       int
	ServerReadTimeout          int
	ServerWriteTimeout         int
	ServerIdleTimeout          int
	ServerShutdownGraceTimeout int
	ServerMaxHeaderBytes       int

	LogLevel    common.LogLevel
	LogFilePath string
	CORSConfig  *common.CORSConfig
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

	env := Environment(common.GetEnvString("APP_ENV", string(EnvDev)))
	if !slices.Contains(Environments[:], env) {
		logger.Panic(fmt.Sprintf("Invalid APP_ENV: %s, supported envs are: %v", env, Environments))
	}

	// server
	host := common.GetEnvString("APP_HOST", "")
	port := common.GetEnvInt("APP_PORT", 8080)

	serverReadTimeout := common.GetEnvInt("APP_SERVER_READ_TIMEOUT_SECONDS", 15)
	serverWriteTimeout := common.GetEnvInt("APP_SERVER_WRITE_TIMEOUT_SECONDS", 15)
	serverIdleTimeout := common.GetEnvInt("APP_SERVER_IDLE_TIMEOUT_SECONDS", 60)
	serverShutdownGraceTimeout := common.GetEnvInt("APP_SERVER_SHUTDOWN_GRACE_TIMEOUT_SECONDS", 10)
	serverMaxHeaderBytes := common.GetEnvInt("APP_SERVER_MAX_HEADER_BYTES", 8192)

	logLevel := common.LogLevel(common.GetEnvString("APP_LOG_LEVEL", string(common.LogLevelInfo)))
	if !slices.Contains(common.LogLevels[:], logLevel) {
		logger.Panic(fmt.Sprintf("Invalid LOG_LEVEL: %s, supported levels are: %v", logLevel, common.LogLevels))
	}
	logFilePath := common.GetEnvString("APP_LOG_FILE", "logs/app.log")

	corsConfig := common.NewCors(logger)

	config = &Config{
		AppName:                    "game-hangar",
		Env:                        env,
		Host:                       host,
		Port:                       port,
		ServerReadTimeout:          serverReadTimeout,
		ServerWriteTimeout:         serverWriteTimeout,
		ServerIdleTimeout:          serverIdleTimeout,
		ServerShutdownGraceTimeout: serverShutdownGraceTimeout,
		ServerMaxHeaderBytes:       serverMaxHeaderBytes,
		LogLevel:                   logLevel,
		LogFilePath:                logFilePath,
		CORSConfig:                 corsConfig,
	}
	return config
}

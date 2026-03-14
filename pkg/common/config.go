package common

import (
	"fmt"
	"slices"

	"github.com/joho/godotenv"
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
	AppName     string
	AppVersion  string
	AppCommitID string
	Env         Environment
	// server configs
	Host                       string
	Port                       int
	ServerReadTimeout          int
	ServerWriteTimeout         int
	ServerIdleTimeout          int
	ServerShutdownGraceTimeout int
	ServerMaxHeaderBytes       int
	CORSConfig                 *CORSConfig

	LogLevel    LogLevel
	LogFilePath string
}

func (c Config) GetAppURL() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c Config) GetVersion() string {
	return fmt.Sprintf("%s+%s", c.AppVersion, c.AppCommitID)
}

func (c Config) GetFullName() string {
	return fmt.Sprintf("%s v%s", c.AppName, c.GetVersion())
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		panic("Config is not initialized. Please initialize the Config before using it.")
	}
	return config
}

func NewConfig(logger Logger, envPrefix string) *Config {
	_ = godotenv.Load(".env")

	appName := GetEnvString(envPrefix+"APP_NAME", "unknown")
	appVersion := GetEnvString(envPrefix+"APP_VERSION", "0.0.1-dev")
	appCommitID := GetEnvString(envPrefix+"APP_COMMIT_ID", "")
	env := Environment(GetEnvString(envPrefix+"APP_ENV", string(EnvDev)))
	if !slices.Contains(Environments[:], env) {
		logger.Panic(
			"Invalid APP_ENV",
			zap.String("env", string(env)),
			zap.Any("supported_envs", Environments),
		)
	}

	// server
	host := GetEnvString(envPrefix+"APP_SERVER_HOST", "")
	port := GetEnvInt(envPrefix+"APP_SERVER_PORT", 8080)
	serverReadTimeout := GetEnvInt(envPrefix+"APP_SERVER_READ_TIMEOUT_SECONDS", 15)
	serverWriteTimeout := GetEnvInt(envPrefix+"APP_SERVER_WRITE_TIMEOUT_SECONDS", 15)
	serverIdleTimeout := GetEnvInt(envPrefix+"APP_SERVER_IDLE_TIMEOUT_SECONDS", 60)
	serverShutdownGraceTimeout := GetEnvInt(envPrefix+"APP_SERVER_SHUTDOWN_GRACE_TIMEOUT_SECONDS", 10)
	serverMaxHeaderBytes := GetEnvInt(envPrefix+"APP_SERVER_MAX_HEADER_BYTES", 8192)

	logLevel := LogLevel(GetEnvString(envPrefix+"APP_LOG_LEVEL", string(LogLevelInfo)))
	if !slices.Contains(LogLevels[:], logLevel) {
		logger.Panic(
			"Invalid LOG_LEVEL",
			zap.String("log_level", string(logLevel)),
			zap.Any("supported_levels", LogLevels),
		)
	}
	logFilePath := GetEnvString(envPrefix+"APP_LOG_FILE", "logs/app.log")

	corsConfig := NewCors(logger)

	config = &Config{
		AppName:                    appName,
		AppVersion:                 appVersion,
		AppCommitID:                appCommitID,
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

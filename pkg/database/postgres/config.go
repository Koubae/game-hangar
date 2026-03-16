package postgres

import (
	"fmt"
	"sync"

	"github.com/koubae/game-hangar/pkg/common"
)

var (
	onceConfig sync.Once
	config     *DatabasePostgresConfig
)

type DatabasePostgresConfig struct {
	Driver                string
	Database              string
	connectionString      string
	host                  string
	port                  int
	user                  string
	password              string
	sslMode               string
	MaxOpenConnections    int32 // Max connections: Use (CPU cores * 2) + 1 as a baseline
	MaxIdleConnections    int32 // Min connections: Keeps warm connections ready for traffic spikes
	MaxConnectionLifetime int32 // MaxConnLifetime: Prevents issues with stale connections or memory leaks
	MaxConnectionIdleTime int32 // MaxConnIdleTime: Closes connections that haven't been used recently
}

func (c *DatabasePostgresConfig) String() string {
	return fmt.Sprintf("DB[%s] database:%s connected @ %s:%d", c.Driver, c.Database, c.host, c.port)

}

// GetConnectionString
// @example
//
//	"postgres://admin:admin@localhost:5432/game_hangar?sslmode=disable"
func (c *DatabasePostgresConfig) GetConnectionString() string {
	if c.connectionString != "" {
		return c.connectionString
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s",
		c.user,
		c.password,
		c.host,
		c.port,
		c.Database,
		c.sslMode,
	)
}

func LoadConfig(envPrefix string) (*DatabasePostgresConfig, error) {
	onceConfig.Do(
		func() {
			config = LoadNewConfig(envPrefix)
		},
	)

	return config, nil
}

func LoadNewConfig(envPrefix string) *DatabasePostgresConfig {
	database := common.GetEnvString(envPrefix+"POSTGRES_DB", "")
	connectionString := common.GetEnvString(envPrefix+"POSTGRES_CONNECTION_STRING", "")
	host := common.GetEnvString(envPrefix+"POSTGRES_HOST", "")
	port := common.GetEnvInt(envPrefix+"POSTGRES_PORT", 0)
	user := common.GetEnvString(envPrefix+"POSTGRES_USER", "")
	password := common.GetEnvString(envPrefix+"POSTGRES_PASS", "")
	sslMode := common.GetEnvString(envPrefix+"POSTGRES_SSL_MODE", "disable")

	// TODO: Test in production these settings are correct. Assumed during development
	maxOpenConnections := common.GetEnvInt(envPrefix+"POSTGRES_MAX_OPEN_CONNECTIONS", 10)
	maxIdleConnections := common.GetEnvInt(envPrefix+"POSTGRES_MAX_IDLE_CONNECTIONS", 2)
	maxConnectionLifetime := common.GetEnvInt(envPrefix+"POSTGRES_MAX_CONNECTION_LIFETIME_MINUTES", 60)
	maxConnectionIdleTime := common.GetEnvInt(envPrefix+"POSTGRES_MAX_CONNECTION_IDLE_TIME_MINUTES", 30)

	return &DatabasePostgresConfig{
		Driver:                "postgres",
		Database:              database,
		connectionString:      connectionString,
		host:                  host,
		port:                  port,
		user:                  user,
		password:              password,
		sslMode:               sslMode,
		MaxOpenConnections:    int32(maxOpenConnections),
		MaxIdleConnections:    int32(maxIdleConnections),
		MaxConnectionLifetime: int32(maxConnectionLifetime),
		MaxConnectionIdleTime: int32(maxConnectionIdleTime),
	}

}

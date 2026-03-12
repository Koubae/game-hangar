package postgres

import (
	"context"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once      sync.Once
	connector *ConnectorPostgres
	errPool   error
)

type ConnectorPostgres struct {
	*pgxpool.Pool
	config *DatabasePostgresConfig
}

func (c *ConnectorPostgres) String() string {
	return c.config.String()
}

func (c *ConnectorPostgres) Ping(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.Pool.Ping(ctx)
}

func (c *ConnectorPostgres) Shutdown() error {
	c.Pool.Close()
	return nil
}

func NewConnector(baseConfig *DatabasePostgresConfig) (*ConnectorPostgres, error) {
	once.Do(
		func() {
			config, err := pgxpool.ParseConfig(baseConfig.GetConnectionString())
			if err != nil {
				errPool = err
				return
			}

			config.MaxConns = baseConfig.MaxOpenConnections
			config.MinConns = baseConfig.MaxIdleConnections
			config.MaxConnLifetime = time.Duration(baseConfig.MaxConnectionLifetime) * time.Minute
			config.MaxConnIdleTime = time.Duration(baseConfig.MaxConnectionIdleTime) * time.Minute
			config.HealthCheckPeriod = 1 * time.Minute         // HealthCheckPeriod: How often the pool checks if connections are still alive
			config.ConnConfig.ConnectTimeout = 5 * time.Second // ConnectTimeout: Time limit for establishing the initial physical connection

			// 3. Create the connection pool
			// Note: NewWithConfig does not immediately connect to the DB
			pool, err := pgxpool.NewWithConfig(context.Background(), config)
			if err != nil {
				errPool = err
				return
			}

			connector = &ConnectorPostgres{
				Pool:   pool,
				config: baseConfig,
			}

			errPool = connector.Ping(context.Background())
			if errPool != nil {
				connector.Close()
				connector = nil
				return
			}
		},
	)
	return connector, errPool
}

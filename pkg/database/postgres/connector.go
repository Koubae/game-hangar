package postgres

import (
	"context"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once    sync.Once
	pool    *pgxpool.Pool
	errPool error
)

func NewConnector(baseConfig *DatabasePostgresConfig) (*pgxpool.Pool, error) {
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
			pool, err = pgxpool.NewWithConfig(context.Background(), config)
			if err != nil {
				pool = nil
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			errPool = pool.Ping(ctx)
			if errPool != nil {
				pool.Close()
				pool = nil
				return
			}
		},
	)
	return pool, errPool
}

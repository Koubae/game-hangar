package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
	ports "github.com/koubae/game-hangar/pkg/database"
)

var (
	once      sync.Once
	connector *ConnectorPostgres
	errPool   error
)

type ConnectorPostgres struct {
	Pool   PoolInterface
	config *DatabasePostgresConfig
}

func InitConnector(baseConfig *DatabasePostgresConfig) (*ConnectorPostgres, error) {
	once.Do(
		func() {
			_connector, err := NewConnector(baseConfig)
			if err != nil {
				errPool = err
				return
			}
			connector = _connector
		},
	)
	return connector, errPool
}

func NewConnector(baseConfig *DatabasePostgresConfig) (*ConnectorPostgres, error) {
	config, err := pgxpool.ParseConfig(baseConfig.GetConnectionString())
	if err != nil {
		return nil, err
	}

	config.MaxConns = baseConfig.MaxOpenConnections
	config.MinConns = baseConfig.MaxIdleConnections
	config.MaxConnLifetime = time.Duration(baseConfig.MaxConnectionLifetime) * time.Minute
	config.MaxConnIdleTime = time.Duration(baseConfig.MaxConnectionIdleTime) * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute         // HealthCheckPeriod: How often the pool checks if connections are still alive
	config.ConnConfig.ConnectTimeout = 5 * time.Second // ConnectTimeout: Time limit for establishing the initial physical connection

	// 3. Create the connection pool
	// Note: NewWithConfig does not immediately connect to the DB
	_pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	_connector := &ConnectorPostgres{
		Pool:   _pool,
		config: baseConfig,
	}

	errPool = _connector.Ping(context.Background())
	if errPool != nil {
		_connector.Shutdown()
		return nil, errPool
	}
	return _connector, nil
}

func (c *ConnectorPostgres) String() string {
	return c.config.String()
}

func (c *ConnectorPostgres) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	return c.Pool.Ping(ctx)
}

func (c *ConnectorPostgres) Shutdown() error {
	c.Pool.Close()
	return nil
}

func (c *ConnectorPostgres) SelectMany(
	ctx context.Context,
	query string,
	args ...any,
) (pgx.Rows, error) {
	return c.Pool.Query(ctx, query, args...)
}

func (c *ConnectorPostgres) SelectOne(
	ctx context.Context,
	query string,
	args ...any,
) pgx.Row {
	return c.Pool.QueryRow(ctx, query, args...)
}

func (c *ConnectorPostgres) SQL(
	ctx context.Context,
	query string,
	args ...any,
) (pgconn.CommandTag, error) {
	return c.Pool.Exec(ctx, query, args...)
}

func (c *ConnectorPostgres) Transaction(
	ctx context.Context,
	txOptions pgx.TxOptions,
) (ports.Transaction, error) {
	tx, err := c.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, err
	}

	return &PostgresTransaction{tx: tx}, nil
}

func (c *ConnectorPostgres) MapDBErrToDomainErr(err error) error {
	return MapPostgresErrToDomainErr(err)
}

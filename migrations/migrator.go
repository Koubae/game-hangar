package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/koubae/game-hangar/pkg/common"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"

	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/koubae/game-hangar/pkg/database/postgres"
)

type Migrator struct {
	db            *sql.DB
	dbPool        *postgres.ConnectorPostgres
	sqlMigrations embed.FS
	Logger        *common.AppLogger
}

func NewMigrator(db *sql.DB, dbPool *postgres.ConnectorPostgres, migrationTable string, sqlMigrations embed.FS, logger *common.AppLogger) *Migrator {
	migrate.SetTable(migrationTable)
	migrate.SetIgnoreUnknown(true)

	return &Migrator{db, dbPool, sqlMigrations, logger}
}

func (m *Migrator) Close() {
	m.dbPool.Shutdown()
	m.db.Close()
}

func (m *Migrator) Run(operation string, limit int) (string, error) {
	switch operation {
	case "up":
		return m.Up(limit)
	case "down":
		return m.Down(limit)
	case "status":
		return m.Status()
	}

	return "MIGRATION_OPERATION_ERROR", fmt.Errorf("invalid operation: %s", operation)
}

func InitializeMigrations(envFile string, appPrefix string, migrationTable string, sqlMigrations embed.FS, createDatabaseFlag bool) *Migrator {
	config := common.NewConfig(common.CreateLogger(common.LogLevelInfo, ""), envFile, appPrefix)
	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	logger.Info("initializing migrations ... for app prefix: ", zap.String("appPrefix", appPrefix))

	dbConfig, err := postgres.LoadConfig(appPrefix)
	if err != nil {
		logger.Fatal("failed to load database configuration", zap.Error(err))
	}

	if createDatabaseFlag {
		createDatabase(appPrefix, dbConfig.Database, logger)
	}

	dbPool, err := postgres.InitConnector(dbConfig)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	} else if dbPool == nil {
		logger.Fatal("Database connection pool is nil...")
	}

	logger.Info("database connection established... ", zap.String("dbConfig", dbConfig.String()))
	return NewMigrator(stdlib.OpenDBFromPool(dbPool.Pool.(*pgxpool.Pool)), dbPool, migrationTable, sqlMigrations, logger)
}

func createDatabase(appPrefix string, database string, logger *common.AppLogger) {
	dbConfigAdmin := postgres.LoadNewConfig(appPrefix)
	dbConfigAdmin.Database = "postgres"

	config, err := pgxpool.ParseConfig(dbConfigAdmin.GetConnectionString())
	if err != nil {
		logger.Fatal("failed to parse admin database configuration", zap.Error(err))
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatal("failed to create admin database connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Info("admin database connection established... ", zap.String("dbConfig", dbConfigAdmin.String()))
	_, err = pool.Exec(context.Background(), "CREATE DATABASE "+database)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "42P04" {
				fmt.Printf("Database %q already exists\n", database)
			}
		} else {
			logger.Fatal("failed to create database %q: %v", zap.String("database", database), zap.Error(err))
		}
	}
	logger.Info("Database ready", zap.String("database", database))

}

func (m *Migrator) Status() (string, error) {
	var migrationsRecords []*migrate.MigrationRecord

	migrationsRecords, err := migrate.GetMigrationRecords(m.db, "postgres")
	if err != nil {
		return "MIGRATION_STATUS_ERROR", err
	}

	for _, migration := range migrationsRecords {
		m.Logger.Info("Migration: ", zap.Any("migration", migration))
	}

	return fmt.Sprintf("MIGRATION_STATUS_OK: %v migrations applied", len(migrationsRecords)), nil
}

func (m *Migrator) Up(limit int) (string, error) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: m.sqlMigrations,
		Root:       "sql",
	}

	appliedMigrations, err := migrate.ExecMax(m.db, "postgres", migrations, migrate.Up, limit)
	if err != nil {
		m.Logger.Error("failed to apply migrations: ", zap.Error(err))
		return "MIGRATION_UP_ERROR", err
	}
	m.Logger.Info("applied migrations: ", zap.Int("appliedMigrations", appliedMigrations))

	return fmt.Sprintf("MIGRATION_UP_OK: %v migrations applied", appliedMigrations), nil
}

func (m *Migrator) Down(limit int) (string, error) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: m.sqlMigrations,
		Root:       "sql",
	}

	appliedRollbacks, err := migrate.ExecMax(m.db, "postgres", migrations, migrate.Down, limit)
	if err != nil {
		m.Logger.Error("failed to apply migrations: ", zap.Error(err))
		return "MIGRATION_DOWN_ERROR", err
	}
	m.Logger.Info("rolled back migrations: ", zap.Int("appliedRollbacks", appliedRollbacks))

	return fmt.Sprintf("MIGRATION_DOWN_OK: %v migrations rolled back", appliedRollbacks), nil
}

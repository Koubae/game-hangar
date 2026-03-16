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
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/koubae/game-hangar/pkg/database/postgres"
)

const (
	MigrationTable = "schema_migrations"
)

type Migrator struct {
	db            *sql.DB
	dbPool        *postgres.ConnectorPostgres
	sqlMigrations embed.FS
	Logger        *common.AppLogger
}

func NewMigrator(db *sql.DB, dbPool *postgres.ConnectorPostgres, sqlMigrations embed.FS, logger *common.AppLogger) *Migrator {
	migrate.SetTable(MigrationTable)
	migrate.SetIgnoreUnknown(true)

	return &Migrator{db, dbPool, sqlMigrations, logger}
}

func (m *Migrator) Close() {
	m.dbPool.Shutdown()
	m.db.Close()
}

func (m *Migrator) Run(operation string) (string, error) {
	switch operation {
	// case "up":
	// 	return Up(m.db)
	// case "down":
	// 	return Down(m.db)
	case "status":
		return m.Status()
	}

	return "MIGRATION_OPERATION_ERROR", fmt.Errorf("invalid operation: %s", operation)
}

// // Up runs all pending migrations.
// func Up(db *sql.DB) (int, error) {
// 	migrations := &migrate.EmbedFileSystemMigrationSource{
// 		FileSystem: sqlMigrations,
// 		Root:       "sql",
// 	}
// 	return migrate.Exec(db, "postgres", migrations, migrate.Up)
// }

// // Down rolls back the last migration.
// func Down(db *sql.DB) (int, error) {
// 	migrations := &migrate.EmbedFileSystemMigrationSource{
// 		FileSystem: sqlMigrations,
// 		Root:       "sql",
// 	}
// 	return migrate.Exec(db, "postgres", migrations, migrate.Down)
// }

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

func InitializeMigrations(appPrefix string, sqlMigrations embed.FS) *Migrator {
	config := common.NewConfig(common.CreateLogger(common.LogLevelInfo, ""), appPrefix)
	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	logger.Info("migrations initialized... for app prefix: ", zap.String("appPrefix", appPrefix))

	dbConfig, err := postgres.LoadConfig(appPrefix)
	if err != nil {
		logger.Fatal("failed to load database configuration", zap.Error(err))
	}
	dbPool, err := postgres.NewConnector(dbConfig)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	} else if dbPool == nil {
		logger.Fatal("Database connection pool is nil...")
	}

	logger.Info("database connection established... ", zap.String("dbConfig", dbConfig.String()))

	_, err = dbPool.Pool.Exec(context.Background(), "CREATE DATABASE "+dbConfig.Database)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "42P04" {
				fmt.Printf("Database %q already exists\n", dbConfig.Database)
			}
		} else {
			logger.Fatal("failed to create database %q: %v", zap.String("database", dbConfig.Database), zap.Error(err))
		}
	}
	logger.Info("Database ready", zap.String("database", dbConfig.Database))

	return NewMigrator(stdlib.OpenDBFromPool(dbPool.Pool), dbPool, sqlMigrations, logger)
}

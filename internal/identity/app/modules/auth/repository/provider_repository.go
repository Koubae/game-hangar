package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/database/postgres"
)

type IProviderRepository interface {
	GetProvider(ctx context.Context, name string) (*model.Provider, error)
}

type ProviderRepository struct {
	DBConnector *postgres.ConnectorPostgres
}

func NewProviderRepository(connector *postgres.ConnectorPostgres) *ProviderRepository {
	return &ProviderRepository{DBConnector: connector}
}

func (r *ProviderRepository) getDB() *sql.DB {
	return stdlib.OpenDBFromPool(r.DBConnector.Pool.(*pgxpool.Pool))
}

func (r *ProviderRepository) GetProvider(ctx context.Context, name string) (*model.Provider, error) {
	db := r.getDB()
	defer db.Close()

	query := "SELECT id, name, display_name, category, disabled, created, updated FROM provider WHERE name = $1"
	row := db.QueryRowContext(ctx, query, name)

	var provider model.Provider
	err := row.Scan(&provider.ID, &provider.Name, &provider.DisplayName, &provider.Category, &provider.Disabled, &provider.Created, &provider.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &provider, nil
}

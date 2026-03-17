package repository

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/koubae/game-hangar/pkg/database/postgres"
)

type IProviderRepository interface {
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

// func (r *ProviderRepository) GetProvider(ctx context.Context, id string) (*model.Provider, error) {
// 	db := r.getDB()
// 	defer db.Close()

// 	query := "SELECT id, name, display_name, category, disabled, created, updated FROM provider WHERE id = $1"
// 	row := db.QueryRowContext(ctx, query, id)

// 	var provider model.Provider
// 	err := row.Scan(&provider.ID, &provider.Name, &provider.DisplayName, &provider.Category, &provider.Disabled, &provider.Created, &provider.Updated)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &provider, nil
// }

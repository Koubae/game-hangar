package repository

import (
	"context"
	"sync"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"go.uber.org/zap"
)

type IProviderRepository interface {
	GetProvider(ctx context.Context, name string) (*model.Provider, error)
}

type ProviderRepository struct {
	DB *postgres.ConnectorPostgres

	mu             sync.RWMutex
	providersCache map[string]*model.Provider
}

func NewProviderRepository(connector *postgres.ConnectorPostgres) *ProviderRepository {
	r := &ProviderRepository{
		DB:             connector,
		providersCache: make(map[string]*model.Provider),
	}
	r.loadProviders(context.Background())
	return r
}

func (r *ProviderRepository) loadProviders(ctx context.Context) {
	logger := common.GetLogger()
	logger.Info("loading providers...")

	query := "SELECT id, name, display_name, category, disabled, created, updated FROM provider"
	rows, err := r.DB.SelectMany(ctx, query)
	if err != nil {
		logger.Error("failed to load providers", zap.Error(err))
		return
	}
	defer rows.Close()

	r.mu.Lock()
	for rows.Next() {
		var p model.Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.DisplayName, &p.Category, &p.Disabled, &p.Created, &p.Updated); err != nil {
			logger.Error("failed to scan provider", zap.Error(err))
			continue
		}
		r.providersCache[p.Name] = &p
	}
	r.mu.Unlock()

	logger.Info("providers loaded", zap.Int("count", len(r.providersCache)))
}

//	func (r *ProviderRepository) getDB() *sql.DB {
//		return stdlib.OpenDBFromPool(r.DB.Pool.(*pgxpool.Pool))
//	}
//
// TODO: 	on commit 5ea82e1 I removed SELECT query here. but i think in case we
//
//					don't hit Cache then:
//	       1. We attempt query
//	       2. If found we add to cache else return nil
func (r *ProviderRepository) GetProvider(ctx context.Context, name string) (*model.Provider, error) {
	r.mu.RLock()
	p, ok := r.providersCache[name]
	r.mu.RUnlock()
	if ok {
		return p, nil
	}
	return nil, nil
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type IProviderRepository interface {
	LoadProviders(ctx context.Context, db database.DBTX)
	GetProvider(ctx context.Context, db database.DBTX, name string) (*model.Provider, error)
}

type ProviderRepository struct {
	mu             sync.RWMutex
	providersCache map[string]*model.Provider
}

func NewProviderRepository() *ProviderRepository {
	r := &ProviderRepository{
		providersCache: make(map[string]*model.Provider),
	}
	return r
}

func (r *ProviderRepository) LoadProviders(ctx context.Context, db database.DBTX) {
	logger := common.GetLogger()
	logger.Info("loading providers...")

	const query = "SELECT id, name, display_name, category, disabled, created, updated FROM provider"

	rows, err := db.SelectMany(ctx, query)
	if err != nil {
		logger.Error("failed to load providers", zap.Error(err))
		return
	}
	defer rows.Close()

	r.mu.Lock()
	for rows.Next() {
		var p model.Provider
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.DisplayName,
			&p.Category,
			&p.Disabled,
			&p.Created,
			&p.Updated,
		); err != nil {
			logger.Error("failed to scan provider", zap.Error(err))
			continue
		}
		r.providersCache[p.Name] = &p
	}
	r.mu.Unlock()

	logger.Info("providers loaded", zap.Int("count", len(r.providersCache)))
}

func (r *ProviderRepository) GetProvider(ctx context.Context, db database.DBTX, name string) (*model.Provider, error) {
	r.mu.RLock()
	m, ok := r.providersCache[name]
	r.mu.RUnlock()
	if ok {
		return m, nil
	}

	logger := common.GetLogger()
	logger.Warn("provider not found in cache, attempt to load from db", zap.String("name", name))

	m, err := r.getProvider(ctx, db, name)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *ProviderRepository) getProvider(ctx context.Context, db database.DBTX, name string) (*model.Provider, error) {
	const query = `
		SELECT id, name, display_name, category, disabled, created, updated 
			FROM provider
		WHERE name = $1 

	`

	var m model.Provider
	if err := db.SelectOne(ctx, query, name).Scan(
		&m.ID,
		&m.Name,
		&m.DisplayName,
		&m.Category,
		&m.Disabled,
		&m.Created,
		&m.Updated,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("error while getProvider, error: %w", err)
	}

	return &m, nil
}

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
	GetProvider(
		ctx context.Context,
		db database.DBTX,
		source string,
		_type string,
	) (*model.Provider, error)
}

type ProviderRepositoryFactory func() IProviderRepository

type ProviderRepository struct {
	mu             sync.RWMutex
	providersCache map[string]map[string]*model.Provider
}

func NewProviderRepository() *ProviderRepository {
	r := &ProviderRepository{
		providersCache: make(map[string]map[string]*model.Provider),
	}
	return r
}

func (r *ProviderRepository) LoadProviders(
	ctx context.Context,
	db database.DBTX,
) {
	logger := common.GetLogger()
	logger.Info("loading providers...")

	const query = "SELECT id, source, type, display_name, category, disabled, created, updated FROM provider"

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
			&p.Source,
			&p.Type,
			&p.DisplayName,
			&p.Category,
			&p.Disabled,
			&p.Created,
			&p.Updated,
		); err != nil {
			logger.Error("failed to scan provider", zap.Error(err))
			continue
		}

		r.addProviderInCache(p.Source, p.Type, &p)

	}
	r.mu.Unlock()

	logger.Info("providers loaded", zap.Int("count", len(r.providersCache)))
}

// addProviderInCache should be called within the r.mu.Lock
func (r *ProviderRepository) addProviderInCache(
	source string,
	_type string,
	p *model.Provider,
) {
	if _, ok := r.providersCache[source]; !ok {
		r.providersCache[source] = make(map[string]*model.Provider)
	}

	r.providersCache[source][p.Type] = p
}

func (r *ProviderRepository) GetProvider(
	ctx context.Context,
	db database.DBTX,
	source string,
	_type string,
) (*model.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.providersCache[source]
	if ok {
		m, ok := s[_type]
		if ok {
			return m, nil
		}

	}

	logger := common.GetLogger()
	logger.Warn(
		"provider not found in cache, attempt to load from db",
		zap.String("source", source),
		zap.String("type", _type),
	)

	m, err := r.getProvider(ctx, db, source, _type)
	if err != nil {
		return nil, err
	}

	r.addProviderInCache(source, _type, m)
	return m, nil
}

func (r *ProviderRepository) getProvider(
	ctx context.Context,
	db database.DBTX,
	source string,
	_type string,
) (*model.Provider, error) {
	const query = `
		SELECT id, source, type, display_name, category, disabled, created, updated 
			FROM provider
		WHERE source = @source AND type = @type 

	`

	var m model.Provider
	if err := db.SelectOne(ctx, query, pgx.StrictNamedArgs{"source": source, "type": _type}).Scan(
		&m.ID,
		&m.Source,
		&m.Type,
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

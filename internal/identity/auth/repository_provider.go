package auth

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"go.uber.org/zap"
)

type IProviderRepository interface {
	LoadProviders(ctx context.Context, db database.DBTX)
	GetProvider(
		ctx context.Context,
		db database.DBTX,
		source string,
		_type string,
	) (*Provider, error)
}

type ProviderRepositoryFactory func() IProviderRepository

type ProviderRepository struct {
	mu             sync.RWMutex
	ProvidersCache map[string]map[string]*Provider
}

func NewProviderRepository() IProviderRepository {
	r := &ProviderRepository{
		ProvidersCache: make(map[string]map[string]*Provider),
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
		logger.DPanic("failed to load providers", zap.Error(err))
		return
	}
	defer rows.Close()

	r.mu.Lock()
	for rows.Next() {
		var p Provider
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

		r.addProviderInCache(&p)

	}
	r.mu.Unlock()

	logger.Info("providers loaded", zap.Int("count", len(r.ProvidersCache)))
}

func (r *ProviderRepository) GetProvider(
	ctx context.Context,
	db database.DBTX,
	source string,
	_type string,
) (*Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.ProvidersCache[source]
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

	r.addProviderInCache(m)
	return m, nil
}

func (r *ProviderRepository) getProvider(
	ctx context.Context,
	db database.DBTX,
	source string,
	_type string,
) (*Provider, error) {
	const query = `
		SELECT id, source, type, display_name, category, disabled, created, updated 
			FROM provider
		WHERE source = @source AND type = @type 

	`

	var m Provider
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
		return nil, errspkg.DBErrToAppErr(db.MapDBErrToDomainErr(err), "provider")
	}

	return &m, nil
}

// addProviderInCache should be called within the r.mu.Lock
func (r *ProviderRepository) addProviderInCache(p *Provider) {
	if _, ok := r.ProvidersCache[p.Source]; !ok {
		r.ProvidersCache[p.Source] = make(map[string]*Provider)
	}

	r.ProvidersCache[p.Source][p.Type] = p
}

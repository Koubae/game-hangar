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

type IPermissionRepository interface {
	LoadPermissions(ctx context.Context, db database.DBTX)
	GetPermissions(
		ctx context.Context,
		db database.DBTX,
		ids []int64,
	) []*Permission
}

type PermissionRepositoryFactory func() IPermissionRepository

type PermissionRepository struct {
	mu               sync.RWMutex
	PermissionsCache map[int64]*Permission
}

func NewPermissionRepository() IPermissionRepository {
	r := &PermissionRepository{
		PermissionsCache: make(map[int64]*Permission),
	}
	return r
}

func (r *PermissionRepository) LoadPermissions(
	ctx context.Context,
	db database.DBTX,
) {
	logger := common.GetLogger()
	logger.Info("loading Permissions...")

	const query = "SELECT id, service, resource, action, created FROM permissions"

	rows, err := db.SelectMany(ctx, query)
	if err != nil {
		logger.DPanic("failed to load Permissions", zap.Error(err))
		return
	}
	defer rows.Close()

	r.mu.Lock()
	for rows.Next() {
		var m Permission
		if err := rows.Scan(
			&m.ID,
			&m.Service,
			&m.Resource,
			&m.Action,
			&m.Created,
		); err != nil {
			logger.Error("failed to scan Permission", zap.Error(err))
			continue
		}

		r.addPermissionInCache(&m)

	}
	r.mu.Unlock()

	logger.Info("Permissions loaded", zap.Int("count", len(r.PermissionsCache)))
}

func (r *PermissionRepository) GetPermissions(
	ctx context.Context,
	db database.DBTX,
	ids []int64,
) []*Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()

	permissions := make([]*Permission, 0)
	for _, id := range ids {
		if s, ok := r.PermissionsCache[id]; ok {
			permissions = append(permissions, s)
			continue
		}

		logger := common.GetLogger()
		logger.Warn(
			"Permission not found in cache, attempt to load from db",
			zap.Int64("id", id),
		)

		m, err := r.getPermission(ctx, db, id)
		if err != nil {
			logger.Error("error loading permission from db", zap.Int64("id", id), zap.Error(err))
			continue
		}

		permissions = append(permissions, m)
		r.addPermissionInCache(m)
	}

	return permissions
}

func (r *PermissionRepository) getPermission(
	ctx context.Context,
	db database.DBTX,
	id int64,
) (*Permission, error) {
	const query = `
		SELECT id, service, resource, action, created
			FROM permissions
		WHERE id = @id
	`

	var m Permission
	if err := db.SelectOne(ctx, query, pgx.StrictNamedArgs{"id": id}).Scan(
		&m.ID,
		&m.Service,
		&m.Resource,
		&m.Action,
		&m.Created,
	); err != nil {
		return nil, errspkg.DBErrToAppErr(db.MapDBErrToDomainErr(err), "Permission")
	}

	return &m, nil
}

// addPermissionInCache should be called within the r.mu.Lock
func (r *PermissionRepository) addPermissionInCache(p *Permission) {
	r.PermissionsCache[p.ID] = p
}

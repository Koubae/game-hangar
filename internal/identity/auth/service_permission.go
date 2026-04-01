package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"go.uber.org/zap"
)

type PermissionService struct {
	db         database.DBTX
	repository IPermissionRepository
}
type PermissionServiceProvider func(d database.DBTX) *PermissionService

type PermissionServiceFactory func(d database.DBTX, r IPermissionRepository) *PermissionService

func NewPermissionService(db database.DBTX, r IPermissionRepository) *PermissionService {
	return &PermissionService{
		db:         db,
		repository: r,
	}
}

func (s *PermissionService) LoadAdminAccountPermissions(ctx context.Context, accountID string, scopeRequested string) (
	authpkg.Permissions,
	string,
	error,
) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger := common.GetLogger()

	permissionsRecords := s.repository.GetAdminAccountPermissions(ctx, s.db, accountID)
	if len(permissionsRecords) == 0 {
		logger.Warn("no permissions found for admin_account", zap.String("accountID", accountID))
		return authpkg.PermissionEmpty, "", errspkg.AuthPermissionsScopeEmpty
	}
	logger.Debug(
		"loaded permissions for admin_account",
		zap.String("accountID", accountID),
		zap.Any("permissions", permissionsRecords),
	)

	scopes := make([]string, len(permissionsRecords))
	for i, permission := range permissionsRecords {
		scopes[i] = fmt.Sprintf("%s:%s:%s", permission.Service, permission.Resource, permission.Action)
	}

	scope := strings.Join(scopes, "|")
	permissions, err := authpkg.ParsePermissions(scope)
	if err != nil {
		logger.Error(
			"failed to parse permissions for admin_account",
			zap.String("accountID", accountID),
			zap.String("scope", scope),
			zap.Error(err),
		)
		return authpkg.PermissionEmpty, "", errspkg.AuthPermissionsScopeEmpty
	}
	return permissions, scope, nil
}

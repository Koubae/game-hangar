package service

import (
	"context"
	"time"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type ProviderService struct {
	db         database.DBTX
	repository repository.IProviderRepository
}

func NewProviderService(d database.DBTX, r repository.IProviderRepository) *ProviderService {
	return &ProviderService{
		db:         d,
		repository: r,
	}
}

func (s *ProviderService) IsProviderEnabled(ctx context.Context, source string, _type string) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger := common.GetLogger()

	provider, err := s.repository.GetProvider(ctx, s.db, source, _type)
	if err != nil {
		logger.Error("error while checking if provider is enabled",
			zap.String("source", source),
			zap.String("type", _type),
			zap.Error(err),
		)
		return false
	}

	return !provider.Disabled
}

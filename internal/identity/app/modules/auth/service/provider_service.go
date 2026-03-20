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
	connector  database.Connector
	repository repository.IProviderRepository
}

func NewProviderService(c database.Connector, r repository.IProviderRepository) *ProviderService {
	return &ProviderService{
		connector:  c,
		repository: r,
	}
}

func (s *ProviderService) IsProviderEnabled(ctx context.Context, name string) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger := common.GetLogger()

	provider, err := s.repository.GetProvider(ctx, s.connector, name)
	if err != nil {
		logger.Error("error while checking if provider is enabled", zap.String("name", name), zap.Error(err))
		return false
	}

	return provider.Disabled == false
}

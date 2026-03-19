package service

import (
	"context"
	"time"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

type ProviderService struct {
	repository repository.IProviderRepository
}

func NewProviderService(r repository.IProviderRepository) *ProviderService {
	return &ProviderService{
		repository: r,
	}
}

func (s *ProviderService) IsProviderEnabled(ctx context.Context, name string) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger := common.GetLogger()

	provider, err := s.repository.GetProvider(ctx, name)
	if err != nil {
		logger.Error("error while checking if provider is enabled", zap.String("name", name), zap.Error(err))
		return false
	}

	return provider.Disabled == false
}

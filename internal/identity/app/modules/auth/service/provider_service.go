package service

import (
	"context"
	"fmt"
	"time"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type ProviderService struct {
	db         database.DBTX
	repository repository.IProviderRepository
}

type ProviderServiceFactory func(d database.DBTX, r repository.IProviderRepository) *ProviderService

func NewProviderService(
	d database.DBTX,
	r repository.IProviderRepository,
) *ProviderService {
	return &ProviderService{
		db:         d,
		repository: r,
	}
}

func (s *ProviderService) IsProviderEnabled(
	ctx context.Context,
	source string,
	_type string,
) bool {
	provider, err := s.GetProvider(ctx, source, _type)
	if err != nil {
		logger := common.GetLogger()
		logger.Error(
			"error while checking if provider is enabled",
			zap.String("source", source),
			zap.String("type", _type),
			zap.Error(err),
		)
		return false
	}

	return !provider.Disabled
}

func (s *ProviderService) GetProvider(
	ctx context.Context,
	source string,
	_type string,
) (*auth.Provider, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	provider, err := s.repository.GetProvider(ctx, s.db, source, _type)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *ProviderService) GetEnabledProvider(
	ctx context.Context,
	source string,
	_type string,
) (*auth.Provider, error) {
	logger := common.GetLogger()
	provider, err := s.GetProvider(
		ctx,
		source,
		_type,
	)
	if err != nil {
		logger.Error(
			"unexpected error while retrieving a provider",
			zap.String(
				"provider",
				fmt.Sprintf("%s.%s", source, _type),
			),
			zap.Error(err),
		)
		return nil, errs.ProviderNotFound

	}
	if provider.Disabled {
		logger.Warn(
			"requested disabled provider source",
			zap.String(
				"provider",
				fmt.Sprintf("%s.%s", source, _type),
			),
		)
		return nil, errs.ProviderDisabled

	}

	return provider, nil
}

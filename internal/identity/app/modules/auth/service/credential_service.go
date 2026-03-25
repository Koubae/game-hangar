package service

import (
	"context"
	"errors"
	"time"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type CredentialService struct {
	db         database.DBTX
	repository repository.ICredentialRepository
}

func NewCredentialService(
	d database.DBTX,
	r repository.ICredentialRepository,
) *CredentialService {
	return &CredentialService{
		db:         d,
		repository: r,
	}
}

func (s *CredentialService) GetCredentialByProvider(
	ctx context.Context,
	providerID int64,
	credential string,
) (*model.AccountCredential, error) {
	return s.getCredentialByProvider(ctx, providerID, credential)
}

func (s *CredentialService) getCredentialByProvider(
	ctx context.Context,
	providerID int64,
	credential string,
) (*model.AccountCredential, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger := common.GetLogger()

	cred, err := s.repository.GetCredentialByProvider(
		ctx,
		s.db,
		providerID,
		credential,
	)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			logger.Error(
				"error while getting credential by provider",
				zap.Int64("providerID", providerID),
				zap.String("credential", credential),
				zap.Error(err),
			)
		}

		return nil, err
	}

	return cred, nil
}

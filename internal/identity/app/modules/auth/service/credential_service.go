package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type CredentialService struct {
	db         database.DBTX
	repository auth.ICredentialRepository
}

type CredentialServiceProvider func(db database.DBTX) *CredentialService

type CredentialServiceFactory func(
	d database.DBTX,
	r auth.ICredentialRepository,
) *CredentialService

func NewCredentialService(
	d database.DBTX,
	r auth.ICredentialRepository,
) *CredentialService {
	return &CredentialService{
		db:         d,
		repository: r,
	}
}

func (s *CredentialService) CreateCredentialTypeUsername(
	ctx context.Context,
	credential string,
	accountID uuid.UUID,
	provider *auth.Provider,
	secret string,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger := common.GetLogger()

	if auth.ProviderType(provider.Type) != auth.Username {
		logger.Warn(
			"[CredentialService] create credential got incorrect provider type",
			zap.String("providerSource", provider.Source),
			zap.String("providerType", provider.Type),
			zap.String("providerName", provider.DisplayName),
		)
		return 0, errs.AccountCredCreateIncorrectProviderType
	}

	verifiedAt := time.Now().UTC()
	params := auth.NewAccountCredential{
		Credential: credential,
		AccountID:  accountID,
		ProviderID: provider.ID,
		Secret:     secret,
		SecretType: "password",
		Verified:   true,
		VerifiedAt: &verifiedAt,
	}

	id, err := s.repository.CreateAccountCredential(ctx, s.db, params)
	if err != nil {
		if !errors.Is(err, errs.ResourceDuplicate) {
			logger.Error(
				"[CredentialService] unexpected error while creating new credential",
				zap.Error(err),
				zap.String("credential", credential),
				zap.String("accountID", accountID.String()),
				zap.String("providerSource", provider.Source),
				zap.String("providerType", provider.Type),
				zap.String("providerName", provider.DisplayName),
			)
		}
		return 0, err
	}

	logger.Debug(
		"[CredentialService] created new credential",
		zap.Int64("id", id),
		zap.String("credential", credential),
		zap.String("accountID", accountID.String()),
		zap.String("providerSource", provider.Source),
		zap.String("providerType", provider.Type),
		zap.String("providerName", provider.DisplayName),
	)
	return id, nil
}

func (s *CredentialService) GetCredentialByProvider(
	ctx context.Context,
	providerID int64,
	credential string,
) (*auth.AccountCredential, error) {
	return s.getCredentialByProvider(ctx, providerID, credential)
}

func (s *CredentialService) getCredentialByProvider(
	ctx context.Context,
	providerID int64,
	credential string,
) (*auth.AccountCredential, error) {
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
		if !errors.Is(err, errs.ResourceNotFound) {
			logger.Error(
				"[CredentialService] error while getting credential by provider",
				zap.Int64("providerID", providerID),
				zap.String("credential", credential),
				zap.Error(err),
			)
		}

		return nil, err
	}

	return cred, nil
}

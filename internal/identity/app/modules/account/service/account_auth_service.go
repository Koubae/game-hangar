package service

import (
	"context"
	"time"

	"github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	authModel "github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	authSrv "github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type AccountAuthService struct {
	db         database.DBTX
	repository repository.IAccountRepository

	providerSrv   *authSrv.ProviderService
	credentialSrv *authSrv.CredentialService
}

type AccountAuthServiceFactory func(
	d database.DBTX,
	r repository.IAccountRepository,
	providerSrv *authSrv.ProviderService,
	credentialSrv *authSrv.CredentialService,
) *AccountAuthService

func NewAccountAuthService(
	d database.DBTX,
	r repository.IAccountRepository,
	providerSrv *authSrv.ProviderService,
	credentialSrv *authSrv.CredentialService,
) *AccountAuthService {
	return &AccountAuthService{
		db:            d,
		repository:    r,
		providerSrv:   providerSrv,
		credentialSrv: credentialSrv,
	}
}

func (s *AccountAuthService) RegisterByUsername(
	ctx context.Context,
	source string,
	credential string,
	secret string,
) error {
	n := "[AccountAuthService.Register]"

	logger := common.GetLogger()
	defer logger.TimeIt("info", n)()
	logger.Info(
		n+" started ...",
		zap.String("source", source),
		zap.String("credential", credential),
	) // TODO: This should be debug + we should mesure stats?

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	provider, err := s.providerSrv.GetEnabledProvider(
		ctx,
		source,
		string(authModel.Username),
	)
	if err != nil {
		return err
	}

	_ = provider

	return nil
}

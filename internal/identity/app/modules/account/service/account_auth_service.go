package service

import (
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	authSrv "github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"
	"github.com/koubae/game-hangar/pkg/database"
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

package container

import (
	authSrv "github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"
	"github.com/koubae/game-hangar/pkg/database"
)

type Scope struct {
	c  *AppContainer
	db database.DBTX
}

func (s Scope) ProviderService() *authSrv.ProviderService {
	return s.c.ProviderService(s.db)
}

func (s Scope) CredentialService() *authSrv.CredentialService {
	return s.c.CredentialService(s.db)
}

package container

import (
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/pkg/database"
)

type Scope struct {
	c  *AppContainer
	db database.DBTX
}

func (s Scope) ProviderService() *auth.ProviderService {
	return s.c.ProviderService(s.db)
}

func (s Scope) CredentialService() *auth.CredentialService {
	return s.c.CredentialService(s.db)
}

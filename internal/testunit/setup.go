package testunit

import (
	"testing"

	identityContainer "github.com/koubae/game-hangar/internal/identity/app/container"
	authSrv "github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/stretchr/testify/require"
)

func Setup() *common.AppLogger {
	logger := common.CreateLogger("dpanic", "/tmp/")
	return logger
}

func NewTestIdentityAppContainer(t *testing.T) *identityContainer.AppContainer {
	logger := Setup()

	connector := MockDBConnector()

	providerRepositoryFactory := NewMockProviderRepository
	credentialRepositoryFactory := NewMockCredentialRepository
	accountRepositoryFactory := NewMockAccountRepository

	// NOTE: Services
	providerServiceFactory := authSrv.NewProviderService
	credentialServiceFactory := authSrv.NewCredentialService

	dependencies := &identityContainer.AppDependencies{
		Logger:    logger,
		Connector: connector,

		ProviderRepositoryFactory:   providerRepositoryFactory,
		CredentialRepositoryFactory: credentialRepositoryFactory,
		AccountRepositoryFactory:    accountRepositoryFactory,

		ProviderServiceFactory:   providerServiceFactory,
		CredentialServiceFactory: credentialServiceFactory,
	}

	container, err := identityContainer.NewAppContainer(
		AppPrefix,
		logger,
		dependencies,
	)
	require.NoError(t, err)

	return container
}

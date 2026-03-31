package testunit

import (
	"net/http"
	"testing"

	auth2 "github.com/koubae/game-hangar/internal/identity/auth"
	identityContainer "github.com/koubae/game-hangar/internal/identity/container"
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
	authServiceFactory := auth2.NewSecretsService
	providerServiceFactory := auth2.NewProviderService
	credentialServiceFactory := auth2.NewCredentialService
	accountAuthServiceFactory := auth2.NewAccountAuthService

	dependencies := &identityContainer.AppDependencies{
		Logger:    logger,
		Connector: connector,

		AuthServiceFactory:          authServiceFactory,
		ProviderRepositoryFactory:   providerRepositoryFactory,
		CredentialRepositoryFactory: credentialRepositoryFactory,
		AccountRepositoryFactory:    accountRepositoryFactory,

		ProviderServiceFactory:    providerServiceFactory,
		CredentialServiceFactory:  credentialServiceFactory,
		AccountAuthServiceFactory: accountAuthServiceFactory,
	}

	container, err := identityContainer.NewAppContainer(
		AppPrefix,
		logger,
		dependencies,
	)
	require.NoError(t, err)

	return container
}

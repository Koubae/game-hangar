package testunit

import (
	"net/http"
	"testing"

	"github.com/koubae/game-hangar/internal/identity"
	auth2 "github.com/koubae/game-hangar/internal/identity/auth"
	identityContainer "github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/web"
	"github.com/stretchr/testify/require"
)

func Setup() (*common.AppLogger, *common.Config) {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	config := common.NewConfig(loggerTmp, EnvFile, AppPrefix)
	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	return logger, config
}

func NewTestIdentityAppContainer(t *testing.T) *identityContainer.AppContainer {
	logger, _ := Setup()

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

func NewTestRouterAndContainer(t *testing.T) (*identityContainer.AppContainer, *http.Handler, *Mocker) {
	_, config := Setup()

	container := NewTestIdentityAppContainer(t)
	handler := web.Router(container, config, identity.RouterRegister(container))

	mocker := NewMocker(container)
	return container, handler, mocker
}

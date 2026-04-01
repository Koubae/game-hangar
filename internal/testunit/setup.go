package testunit

import (
	"net/http"
	"testing"
	"time"

	"github.com/koubae/game-hangar/internal/identity"
	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/identity/auth"
	identityContainer "github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/web"
	"github.com/stretchr/testify/require"
)

const (
	AuthTokenExpirationTime = time.Hour * 4
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
	permissionRepositoryFactory := NewMockPermissionRepository
	accountRepositoryFactory := NewMockAccountRepository

	// NOTE: Services
	authServiceFactory := auth.NewSecretsService
	providerServiceFactory := auth.NewProviderService
	credentialServiceFactory := auth.NewCredentialService
	accountAuthServiceFactory := auth.NewAccountAuthService
	accountManagementServiceFactory := account.NewManagementService

	dependencies := &identityContainer.AppDependencies{
		Logger:    logger,
		Connector: connector,

		AuthServiceFactory:          authServiceFactory,
		ProviderRepositoryFactory:   providerRepositoryFactory,
		CredentialRepositoryFactory: credentialRepositoryFactory,
		PermissionRepositoryFactory: permissionRepositoryFactory,
		AccountRepositoryFactory:    accountRepositoryFactory,

		ProviderServiceFactory:          providerServiceFactory,
		CredentialServiceFactory:        credentialServiceFactory,
		AccountAuthServiceFactory:       accountAuthServiceFactory,
		AccountManagementServiceFactory: accountManagementServiceFactory,
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

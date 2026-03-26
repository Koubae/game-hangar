package testunit

import (
	"testing"

	identityContainer "github.com/koubae/game-hangar/internal/identity/app/container"
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

	dependencies := &identityContainer.AppDependencies{
		Logger:    logger,
		Connector: connector,

		ProviderRepositoryFactory: providerRepositoryFactory,
	}

	container, err := identityContainer.NewAppContainer(
		AppPrefix,
		logger,
		dependencies,
	)
	require.NoError(t, err)

	return container
}

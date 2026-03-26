package container

import (
	"context"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityAuthContainer interface {
	// ProviderService(db database.DBTX) authRepository.IProviderRepository
	ProviderRepository() repository.IProviderRepository
	CredentialRepository() repository.ICredentialRepository
}

type IdentityContainer interface {
	di.Container
	IdentityAuthContainer

	WithDB(db database.DBTX) Scope
	DB() *postgres.ConnectorPostgres

	// Repositories
	// ProviderRepository() repository.IProviderRepository
}

type AppContainer struct {
	logger    common.Logger
	connector *postgres.ConnectorPostgres

	// Repositories
	providerRepository        repository.IProviderRepository
	providerRepositoryFactory func() repository.IProviderRepository

	credentialRepository        repository.ICredentialRepository
	credentialRepositoryFactory func() repository.ICredentialRepository
}

func NewAppContainer(
	appPrefix string,
	logger common.Logger,
) (*AppContainer, error) {
	dbConfig, err := postgres.LoadConfig(appPrefix)
	if err != nil {
		return nil, err
	}
	connector, err := postgres.InitConnector(dbConfig)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"database connection established",
		zap.String("db", connector.String()),
	)

	providerRepositoryFactory := func() repository.IProviderRepository {
		return repository.NewProviderRepository()
	}
	providerRepository := providerRepositoryFactory()
	providerRepository.LoadProviders(context.TODO(), connector)

	credentialRepositoryFactory := func() repository.ICredentialRepository {
		return repository.NewCredentialRepository()
	}

	return &AppContainer{
		logger:                      logger,
		connector:                   connector,
		providerRepository:          providerRepository,
		credentialRepositoryFactory: credentialRepositoryFactory,
	}, nil
}

// NOTE:
// ------------------------------------------
//	Implements di.Container interface
// ------------------------------------------

func (c *AppContainer) Logger() common.Logger {
	return c.logger
}

func (c *AppContainer) Shutdown() error {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	defer func() {
		if z, ok := c.Logger().(*common.AppLogger); ok {
			z.LogCloser(loggerTmp, z.Logger)
		}
	}()

	c.DB().Shutdown()
	c.Logger().Info("database connection closed")
	return nil
}

// NOTE:
// ------------------------------------------
// 	Implements IdentityContainer interface
// ------------------------------------------

func (c *AppContainer) WithDB(db database.DBTX) Scope {
	if db == nil {
		db = c.connector
	}
	return Scope{c: c, db: db}
}

func (c *AppContainer) DB() *postgres.ConnectorPostgres {
	return c.connector
}

// NOTE:
// ------------------------------------------
// 	Dependencies Provider's Factories
// ------------------------------------------

func (c *AppContainer) ProviderRepository() repository.IProviderRepository {
	return c.providerRepository
}

func (c *AppContainer) CredentialRepository() repository.ICredentialRepository {
	if c.credentialRepository == nil {
		c.credentialRepository = c.credentialRepositoryFactory()
	}

	return c.credentialRepository
}

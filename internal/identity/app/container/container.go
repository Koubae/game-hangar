package container

import (
	"context"

	accountRepo "github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	authRepo "github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityAuthContainer interface {
	// ProviderService(db database.DBTX) authRepository.IProviderRepository
	ProviderRepository() authRepo.IProviderRepository
	CredentialRepository() authRepo.ICredentialRepository
	AccountRepository() accountRepo.IAccountRepository
}

type IdentityContainer interface {
	di.Container
	IdentityAuthContainer

	WithDB(db database.DBTX) Scope
	DB() *postgres.ConnectorPostgres
}

type AppContainer struct {
	logger    common.Logger
	connector *postgres.ConnectorPostgres

	// NOTE: Repositories
	providerRepository        authRepo.IProviderRepository
	providerRepositoryFactory authRepo.ProviderRepositoryFactory

	credentialRepository        authRepo.ICredentialRepository
	credentialRepositoryFactory authRepo.CredentialRepositoryFactory

	accountRepository        accountRepo.IAccountRepository
	accountRepositoryFactory accountRepo.AccountRepositoryFactory
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

	providerRepositoryFactory := func() authRepo.IProviderRepository {
		return authRepo.NewProviderRepository()
	}
	providerRepository := providerRepositoryFactory()
	providerRepository.LoadProviders(context.TODO(), connector)

	credentialRepositoryFactory := func() authRepo.ICredentialRepository {
		return authRepo.NewCredentialRepository()
	}

	accountRepositoryFactory := func() accountRepo.IAccountRepository {
		return accountRepo.NewAccountRepository()
	}

	return &AppContainer{
		logger:                      logger,
		connector:                   connector,
		providerRepository:          providerRepository,
		credentialRepositoryFactory: credentialRepositoryFactory,
		accountRepositoryFactory:    accountRepositoryFactory,
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

func (c *AppContainer) ProviderRepository() authRepo.IProviderRepository {
	return c.providerRepository
}

func (c *AppContainer) CredentialRepository() authRepo.ICredentialRepository {
	if c.credentialRepository == nil {
		c.credentialRepository = c.credentialRepositoryFactory()
	}

	return c.credentialRepository
}

func (c *AppContainer) AccountRepository() accountRepo.IAccountRepository {
	if c.accountRepository == nil {
		c.accountRepository = c.accountRepositoryFactory()
	}

	return c.accountRepository
}

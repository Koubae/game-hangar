package container

import (
	accountRepo "github.com/koubae/game-hangar/internal/identity/account"
	auth2 "github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityAuthContainer interface {
	ProviderRepository() auth2.IProviderRepository
	CredentialRepository() auth2.ICredentialRepository

	SecretsService() *auth2.SecretsService
	ProviderService(db database.DBTX) *auth2.ProviderService
	CredentialService(db database.DBTX) *auth2.CredentialService
}

type IdentityAccountContainer interface {
	AccountRepository() accountRepo.IAccountRepository
	AccountAuthService(db database.Connector) *auth2.AccountAuthService
}

type IdentityContainer interface {
	di.Container
	IdentityAuthContainer
	IdentityAccountContainer

	WithDB(db database.DBTX) Scope
	DB() *postgres.ConnectorPostgres
}

type AppContainer struct {
	logger    common.Logger
	connector *postgres.ConnectorPostgres

	// NOTE: Repositories
	providerRepository        auth2.IProviderRepository
	providerRepositoryFactory auth2.ProviderRepositoryFactory

	credentialRepository        auth2.ICredentialRepository
	credentialRepositoryFactory auth2.CredentialRepositoryFactory

	accountRepository        accountRepo.IAccountRepository
	accountRepositoryFactory accountRepo.AccountRepositoryFactory

	// NOTE: Services
	authService               *auth2.SecretsService
	authServiceFactory        auth2.SecretsServiceFactory
	providerServiceFactory    auth2.ProviderServiceFactory
	credentialServiceFactory  auth2.CredentialServiceFactory
	accountAuthServiceFactory auth2.AccountAuthServiceFactory
}

type AppDependencies struct {
	Logger    common.Logger
	Connector *postgres.ConnectorPostgres

	ProviderRepositoryFactory   auth2.ProviderRepositoryFactory
	CredentialRepositoryFactory auth2.CredentialRepositoryFactory
	AccountRepositoryFactory    accountRepo.AccountRepositoryFactory
	AuthServiceFactory          auth2.SecretsServiceFactory
	ProviderServiceFactory      auth2.ProviderServiceFactory
	CredentialServiceFactory    auth2.CredentialServiceFactory
	AccountAuthServiceFactory   auth2.AccountAuthServiceFactory
}

func NewAppContainer(
	appPrefix string,
	logger common.Logger,
	dependencies *AppDependencies,
) (*AppContainer, error) {
	var err error

	if dependencies == nil {
		dependencies, err = createProductionAppDependencies(appPrefix, logger)
		if err != nil {
			return nil, err
		}
	}

	providerRepository := dependencies.ProviderRepositoryFactory()

	return &AppContainer{
		logger:                      dependencies.Logger,
		connector:                   dependencies.Connector,
		providerRepository:          providerRepository,
		providerRepositoryFactory:   dependencies.ProviderRepositoryFactory,
		credentialRepositoryFactory: dependencies.CredentialRepositoryFactory,
		accountRepositoryFactory:    dependencies.AccountRepositoryFactory,
		authServiceFactory:          dependencies.AuthServiceFactory,
		providerServiceFactory:      dependencies.ProviderServiceFactory,
		credentialServiceFactory:    dependencies.CredentialServiceFactory,
		accountAuthServiceFactory:   dependencies.AccountAuthServiceFactory,
	}, nil
}

func createProductionAppDependencies(
	appPrefix string,
	logger common.Logger,
) (*AppDependencies, error) {
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

	return LoadAppDependenciesWithDefaFactories(logger, connector)
}

func LoadAppDependenciesWithDefaFactories(
	logger common.Logger,
	connector *postgres.ConnectorPostgres,
) (*AppDependencies, error) {
	// NOTE: Repositories
	providerRepositoryFactory := auth2.NewProviderRepository
	credentialRepositoryFactory := auth2.NewCredentialRepository
	accountRepositoryFactory := accountRepo.NewAccountRepository

	// NOTE: Services
	authServiceFactory := auth2.NewSecretsService
	providerServiceFactory := auth2.NewProviderService
	credentialServiceFactory := auth2.NewCredentialService
	accountAuthServiceFactory := auth2.NewAccountAuthService

	return &AppDependencies{
		Logger:    logger,
		Connector: connector,

		ProviderRepositoryFactory:   providerRepositoryFactory,
		CredentialRepositoryFactory: credentialRepositoryFactory,
		AccountRepositoryFactory:    accountRepositoryFactory,

		AuthServiceFactory:        authServiceFactory,
		ProviderServiceFactory:    providerServiceFactory,
		CredentialServiceFactory:  credentialServiceFactory,
		AccountAuthServiceFactory: accountAuthServiceFactory,
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

func (c *AppContainer) ProviderRepository() auth2.IProviderRepository {
	return c.providerRepository
}

func (c *AppContainer) CredentialRepository() auth2.ICredentialRepository {
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

func (c *AppContainer) SecretsService() *auth2.SecretsService {
	if c.authService == nil {
		c.authService = c.authServiceFactory()
	}
	return c.authService
}

func (c *AppContainer) ProviderService(
	db database.DBTX,
) *auth2.ProviderService {
	if db == nil {
		db = c.connector
	}
	return c.providerServiceFactory(db, c.ProviderRepository())
}

func (c *AppContainer) CredentialService(
	db database.DBTX,
) *auth2.CredentialService {
	if db == nil {
		db = c.connector
	}
	return c.credentialServiceFactory(db, c.CredentialRepository())
}

func (c *AppContainer) AccountAuthService(
	db database.Connector,
) *auth2.AccountAuthService {
	if db == nil {
		db = c.connector
	}

	providerSrv := c.ProviderService(db)
	repository := c.AccountRepository()

	return c.accountAuthServiceFactory(
		db,
		repository,
		providerSrv,
		func(db database.DBTX) *auth2.CredentialService {
			return c.CredentialService(db)
		},
	)
}

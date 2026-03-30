package container

import (
	accountRepo "github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityAuthContainer interface {
	ProviderRepository() auth.IProviderRepository
	CredentialRepository() auth.ICredentialRepository

	AuthService() *auth.SecretsService
	ProviderService(db database.DBTX) *auth.ProviderService
	CredentialService(db database.DBTX) *auth.CredentialService
}

type IdentityAccountContainer interface {
	AccountRepository() accountRepo.IAccountRepository
	AccountAuthService(db database.Connector) *auth.AccountAuthService
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
	providerRepository        auth.IProviderRepository
	providerRepositoryFactory auth.ProviderRepositoryFactory

	credentialRepository        auth.ICredentialRepository
	credentialRepositoryFactory auth.CredentialRepositoryFactory

	accountRepository        accountRepo.IAccountRepository
	accountRepositoryFactory accountRepo.AccountRepositoryFactory

	// NOTE: Services
	authService               *auth.SecretsService
	authServiceFactory        auth.SecretsServiceFactory
	providerServiceFactory    auth.ProviderServiceFactory
	credentialServiceFactory  auth.CredentialServiceFactory
	accountAuthServiceFactory auth.AccountAuthServiceFactory
}

type AppDependencies struct {
	Logger    common.Logger
	Connector *postgres.ConnectorPostgres

	ProviderRepositoryFactory   auth.ProviderRepositoryFactory
	CredentialRepositoryFactory auth.CredentialRepositoryFactory
	AccountRepositoryFactory    accountRepo.AccountRepositoryFactory
	AuthServiceFactory          auth.SecretsServiceFactory
	ProviderServiceFactory      auth.ProviderServiceFactory
	CredentialServiceFactory    auth.CredentialServiceFactory
	AccountAuthServiceFactory   auth.AccountAuthServiceFactory
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
	providerRepositoryFactory := auth.NewProviderRepository
	credentialRepositoryFactory := auth.NewCredentialRepository
	accountRepositoryFactory := accountRepo.NewAccountRepository

	// NOTE: Services
	authServiceFactory := auth.NewSecretsService
	providerServiceFactory := auth.NewProviderService
	credentialServiceFactory := auth.NewCredentialService
	accountAuthServiceFactory := auth.NewAccountAuthService

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

func (c *AppContainer) ProviderRepository() auth.IProviderRepository {
	return c.providerRepository
}

func (c *AppContainer) CredentialRepository() auth.ICredentialRepository {
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

func (c *AppContainer) AuthService() *auth.SecretsService {
	if c.authService == nil {
		c.authService = c.authServiceFactory()
	}
	return c.authService
}

func (c *AppContainer) ProviderService(
	db database.DBTX,
) *auth.ProviderService {
	if db == nil {
		db = c.connector
	}
	return c.providerServiceFactory(db, c.ProviderRepository())
}

func (c *AppContainer) CredentialService(
	db database.DBTX,
) *auth.CredentialService {
	if db == nil {
		db = c.connector
	}
	return c.credentialServiceFactory(db, c.CredentialRepository())
}

func (c *AppContainer) AccountAuthService(
	db database.Connector,
) *auth.AccountAuthService {
	if db == nil {
		db = c.connector
	}

	providerSrv := c.ProviderService(db)
	repository := c.AccountRepository()

	return c.accountAuthServiceFactory(
		db,
		repository,
		providerSrv,
		func(db database.DBTX) *auth.CredentialService {
			return c.CredentialService(db)
		},
	)
}

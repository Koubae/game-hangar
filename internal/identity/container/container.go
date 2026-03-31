package container

import (
	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityAuthContainer interface {
	ProviderRepository() auth.IProviderRepository
	CredentialRepository() auth.ICredentialRepository

	SecretsService() *auth.SecretsService
	ProviderService(db database.DBTX) *auth.ProviderService
	CredentialService(db database.DBTX) *auth.CredentialService
	AccountAuthService(db database.Connector) *auth.AccountAuthService
}

type IdentityAccountContainer interface {
	AccountRepository() account.IAccountRepository
	AccountManagementService(db database.DBTX) *account.ManagementService
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

	accountRepository        account.IAccountRepository
	accountRepositoryFactory account.AccountRepositoryFactory

	// NOTE: Services
	authService                     *auth.SecretsService
	authServiceFactory              auth.SecretsServiceFactory
	providerServiceFactory          auth.ProviderServiceFactory
	credentialServiceFactory        auth.CredentialServiceFactory
	accountAuthServiceFactory       auth.AccountAuthServiceFactory
	accountManagementServiceFactory account.ManagementServiceFactory
}

type AppDependencies struct {
	Logger    common.Logger
	Connector *postgres.ConnectorPostgres

	ProviderRepositoryFactory       auth.ProviderRepositoryFactory
	CredentialRepositoryFactory     auth.CredentialRepositoryFactory
	AccountRepositoryFactory        account.AccountRepositoryFactory
	AuthServiceFactory              auth.SecretsServiceFactory
	ProviderServiceFactory          auth.ProviderServiceFactory
	CredentialServiceFactory        auth.CredentialServiceFactory
	AccountAuthServiceFactory       auth.AccountAuthServiceFactory
	AccountManagementServiceFactory account.ManagementServiceFactory
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

	auth.LoadPasswordRulesConfig(appPrefix)
	err = authpkg.LoadCerts(appPrefix)
	if err != nil {
		return nil, err
	}

	return &AppContainer{
		logger:                          dependencies.Logger,
		connector:                       dependencies.Connector,
		providerRepository:              providerRepository,
		providerRepositoryFactory:       dependencies.ProviderRepositoryFactory,
		credentialRepositoryFactory:     dependencies.CredentialRepositoryFactory,
		accountRepositoryFactory:        dependencies.AccountRepositoryFactory,
		authServiceFactory:              dependencies.AuthServiceFactory,
		providerServiceFactory:          dependencies.ProviderServiceFactory,
		credentialServiceFactory:        dependencies.CredentialServiceFactory,
		accountAuthServiceFactory:       dependencies.AccountAuthServiceFactory,
		accountManagementServiceFactory: dependencies.AccountManagementServiceFactory,
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

	return LoadAppDependenciesWithDefaultFactories(logger, connector)
}

func LoadAppDependenciesWithDefaultFactories(
	logger common.Logger,
	connector *postgres.ConnectorPostgres,
) (*AppDependencies, error) {
	// NOTE: Repositories
	providerRepositoryFactory := auth.NewProviderRepository
	credentialRepositoryFactory := auth.NewCredentialRepository
	accountRepositoryFactory := account.NewAccountRepository

	// NOTE: Services
	authServiceFactory := auth.NewSecretsService
	providerServiceFactory := auth.NewProviderService
	credentialServiceFactory := auth.NewCredentialService
	accountAuthServiceFactory := auth.NewAccountAuthService
	accountManagementServiceFactory := account.NewManagementService

	return &AppDependencies{
		Logger:    logger,
		Connector: connector,

		ProviderRepositoryFactory:   providerRepositoryFactory,
		CredentialRepositoryFactory: credentialRepositoryFactory,
		AccountRepositoryFactory:    accountRepositoryFactory,

		AuthServiceFactory:              authServiceFactory,
		ProviderServiceFactory:          providerServiceFactory,
		CredentialServiceFactory:        credentialServiceFactory,
		AccountAuthServiceFactory:       accountAuthServiceFactory,
		AccountManagementServiceFactory: accountManagementServiceFactory,
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

	_ = c.DB().Shutdown()
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

func (c *AppContainer) AccountRepository() account.IAccountRepository {
	if c.accountRepository == nil {
		c.accountRepository = c.accountRepositoryFactory()
	}

	return c.accountRepository
}

func (c *AppContainer) SecretsService() *auth.SecretsService {
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

func (c *AppContainer) AccountManagementService(
	db database.DBTX,
) *account.ManagementService {
	if db == nil {
		db = c.connector
	}
	return c.accountManagementServiceFactory(db, c.AccountRepository())
}

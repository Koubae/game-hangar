package container

import (
	accountRepo "github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	authRepo "github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	authSrv "github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/pkg/database/postgres"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type IdentityAuthContainer interface {
	ProviderRepository() authRepo.IProviderRepository
	CredentialRepository() authRepo.ICredentialRepository

	ProviderService(db database.DBTX) *authSrv.ProviderService
	CredentialService(db database.DBTX) *authSrv.CredentialService
}

type IdentityAccountContainer interface {
	AccountRepository() accountRepo.IAccountRepository
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
	providerRepository        authRepo.IProviderRepository
	providerRepositoryFactory authRepo.ProviderRepositoryFactory

	credentialRepository        authRepo.ICredentialRepository
	credentialRepositoryFactory authRepo.CredentialRepositoryFactory

	accountRepository        accountRepo.IAccountRepository
	accountRepositoryFactory accountRepo.AccountRepositoryFactory

	// NOTE: Services
	providerServiceFactory   authSrv.ProviderServiceFactory
	credentialServiceFactory authSrv.CredentialServiceFactory
}

type AppDependencies struct {
	Logger    common.Logger
	Connector *postgres.ConnectorPostgres

	ProviderRepositoryFactory   authRepo.ProviderRepositoryFactory
	CredentialRepositoryFactory authRepo.CredentialRepositoryFactory
	AccountRepositoryFactory    accountRepo.AccountRepositoryFactory
	ProviderServiceFactory      authSrv.ProviderServiceFactory
	CredentialServiceFactory    authSrv.CredentialServiceFactory
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
		providerServiceFactory:      dependencies.ProviderServiceFactory,
		credentialServiceFactory:    dependencies.CredentialServiceFactory,
	}, nil
}

func createProductionAppDependencies(appPrefix string,
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

	// NOTE: Repositories
	providerRepositoryFactory := authRepo.NewProviderRepository
	credentialRepositoryFactory := authRepo.NewCredentialRepository
	accountRepositoryFactory := accountRepo.NewAccountRepository

	// NOTE: Services
	providerServiceFactory := authSrv.NewProviderService
	credentialServiceFactory := authSrv.NewCredentialService

	return &AppDependencies{
		Logger:    logger,
		Connector: connector,

		ProviderRepositoryFactory:   providerRepositoryFactory,
		CredentialRepositoryFactory: credentialRepositoryFactory,
		AccountRepositoryFactory:    accountRepositoryFactory,

		ProviderServiceFactory:   providerServiceFactory,
		CredentialServiceFactory: credentialServiceFactory,
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

func (c *AppContainer) ProviderService(
	db database.DBTX,
) *authSrv.ProviderService {
	if db == nil {
		db = c.connector
	}
	return c.providerServiceFactory(db, c.ProviderRepository())
}

func (c *AppContainer) CredentialService(
	db database.DBTX,
) *authSrv.CredentialService {
	if db == nil {
		db = c.connector
	}
	return c.credentialServiceFactory(db, c.CredentialRepository())
}

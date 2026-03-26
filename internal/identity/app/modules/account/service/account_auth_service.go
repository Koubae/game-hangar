package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/repository"
	authModel "github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	authSrv "github.com/koubae/game-hangar/internal/identity/app/modules/auth/service"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type AccountAuthService struct {
	db         database.Connector
	repository repository.IAccountRepository

	providerSrv           *authSrv.ProviderService
	credentialSrvProvider authSrv.CredentialServiceProvider
}

type AccountAuthServiceFactory func(
	d database.Connector,
	r repository.IAccountRepository,
	providerSrv *authSrv.ProviderService,
	credentialSrvProvider authSrv.CredentialServiceProvider,
) *AccountAuthService

func NewAccountAuthService(
	d database.Connector,
	r repository.IAccountRepository,
	providerSrv *authSrv.ProviderService,

	credentialSrvProvider authSrv.CredentialServiceProvider,
) *AccountAuthService {
	return &AccountAuthService{
		db:                    d,
		repository:            r,
		providerSrv:           providerSrv,
		credentialSrvProvider: credentialSrvProvider,
	}
}

var ErrRegistrationCredExists = errors.New(
	"registration error: credential already exists",
)

func (s *AccountAuthService) RegisterByUsername(
	ctx context.Context,
	source string,
	credential string,
	secret string,
) error {
	n := "[AccountAuthService.RegisterByUsername]"

	logger := common.GetLogger()
	defer logger.TimeIt("info", n)()
	logger.Info(
		n+" started ...",
		zap.String("source", source),
		zap.String("credential", credential),
	) // TODO: This should be debug + we should mesure stats?

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	provider, err := s.providerSrv.GetEnabledProvider(
		ctx,
		source,
		string(authModel.Username),
	)
	if err != nil {
		return err
	}

	cred, err := s.credentialSrvProvider(s.db).GetCredentialByProvider(
		ctx,
		provider.ID,
		credential,
	)
	if cred != nil {
		logger.Warn(n+"attempt to create account using existing credentials",
			zap.String("source", source),
			zap.String("credential", credential),
		)
		return ErrRegistrationCredExists
	} else if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return err
		}
	}

	tx, err := s.db.Transaction(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return &database.ErrOpenTransaction{Err: err}
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil &&
			!errors.Is(rbErr, pgx.ErrTxClosed) {
			logger.Error(n+" error on TX Rollback", zap.Error(rbErr))
		}
	}()

	// NOTE: TRANSACTION BEGIN
	// WARN: All below operation are within a transaction

	if err = tx.Commit(ctx); err != nil {
		logger.Error(
			n+" error on commit",
			zap.Bool(
				"isRollbackError",
				errors.Is(err, pgx.ErrTxCommitRollback),
			),
			zap.Error(err),
		)
		return err
	}

	return nil
}

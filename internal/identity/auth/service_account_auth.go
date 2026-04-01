package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type AccountAuthService struct {
	db         database.Connector
	repository account.IAccountRepository

	providerSrv           *ProviderService
	credentialSrvProvider CredentialServiceProvider
}

type AccountAuthServiceFactory func(
	d database.Connector,
	r account.IAccountRepository,
	providerSrv *ProviderService,
	credentialSrvProvider CredentialServiceProvider,
) *AccountAuthService

func NewAccountAuthService(
	d database.Connector,
	r account.IAccountRepository,
	providerSrv *ProviderService,

	credentialSrvProvider CredentialServiceProvider,
) *AccountAuthService {
	return &AccountAuthService{
		db:                    d,
		repository:            r,
		providerSrv:           providerSrv,
		credentialSrvProvider: credentialSrvProvider,
	}
}

func (s *AccountAuthService) RegisterByUsername(
	ctx context.Context,
	source string,
	credential string,
	secret string,
) (*string, *int64, error) {
	n := "[AccountAuthService.RegisterByUsername] "

	logger := common.GetLogger()
	defer logger.TimeIt("info", n)()
	logger.Info(
		n+"started ...",
		zap.String("source", source),
		zap.String("credential", credential),
	)

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	createAccountParams := account.NewAccount{
		Username: credential,
		Email:    nil,
	}
	if err := createAccountParams.Validate(); err != nil {
		return nil, nil, err
	}

	provider, err := s.providerSrv.GetEnabledProvider(
		ctx,
		source,
		string(Username),
	)
	if err != nil {
		return nil, nil, err
	}

	cred, err := s.credentialSrvProvider(s.db).GetCredentialByProvider(
		ctx,
		provider.ID,
		credential,
	)
	if cred != nil {
		logger.Warn(
			n+"attempt to create account using existing credentials",
			zap.String("source", source),
			zap.String("credential", credential),
		)
		return nil, nil, errs.AccountCredDuplicate
	} else if err != nil {
		if !errors.Is(err, errs.ResourceNotFound) {
			return nil, nil, err
		}
	}

	tx, err := s.db.Transaction(
		ctx, pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil &&
			!errors.Is(rbErr, pgx.ErrTxClosed) {
			logger.Error(n+"error on TX Rollback", zap.Error(rbErr))
		}
	}()

	// NOTE: TRANSACTION BEGIN
	// WARN: All below operation are within a transaction
	credServiceTX := s.credentialSrvProvider(tx)

	id, err := s.repository.CreateAccount(ctx, tx, createAccountParams)
	if err != nil {
		lvl := "error"
		if errs.IsAny(err, errs.ResourceDuplicate) {
			lvl = "debug"
		}

		logger.L(lvl, n+"could not create account, rolling back account", zap.Error(err))
		return nil, nil, err
	}

	accountID, _ := uuid.Parse(*id)
	credID, err := credServiceTX.CreateCredentialTypeUsername(
		ctx,
		credential,
		accountID,
		provider,
		secret,
	)
	if err != nil {
		lvl := "error"
		if errs.IsAny(err, errs.AccountCredCreateIncorrectProviderType, errs.ResourceDuplicate) {
			lvl = "debug"
		}

		logger.L(
			lvl, n+"error while creating credential, rolling back account",
			zap.String("accountID", *id),
			zap.Error(err),
		)
		return nil, nil, err

	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error(
			n+"error on commit",
			zap.Bool(
				"isRollbackError",
				errors.Is(err, pgx.ErrTxCommitRollback),
			),
			zap.Error(err),
		)
		return nil, nil, err
	}

	logger.Info(
		"created new account using username credentials",
		zap.String("accountID", *id),
		zap.Int64("credentialID", credID),
		zap.String("source", source),
		zap.String("credential", credential),
	)
	return id, &credID, nil
}

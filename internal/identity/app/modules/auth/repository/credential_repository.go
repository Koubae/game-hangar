package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/database"
)

type ICredentialRepository interface {
	CreateAccountCredential(
		ctx context.Context,
		db database.DBTX,
		params NewAccountCredential,
	) (int64, error)
	GetCredentialByProvider(
		ctx context.Context,
		db database.DBTX,
		providerID int64,
		credential string,
	) (*model.AccountCredential, error)
}

type CredentialRepositoryFactory func() ICredentialRepository

type NewAccountCredential struct {
	Credential string
	AccountID  uuid.UUID
	ProviderID int64
	Secret     string
	SecretType string
	Verified   bool
	VerifiedAt *time.Time
}

func (p *NewAccountCredential) Validate() error {
	if p.Verified && p.VerifiedAt == nil {
		return &errs.AppError{
			Err: errs.AccountCredVerifiedAtRequired,
			Msg: "verified_at is required when verified is true",
		}
	}
	if !p.Verified && p.VerifiedAt != nil {
		return &errs.AppError{
			Err: errs.AccountCredVerifiedNilWhenIsFalse,
			Msg: "verified_at must be nil when verified is false",
		}
	}
	return nil
}

type CredentialRepository struct{}

func NewCredentialRepository() ICredentialRepository {
	r := &CredentialRepository{}
	return r
}

func (r *CredentialRepository) CreateAccountCredential(
	ctx context.Context,
	db database.DBTX,
	params NewAccountCredential,
) (int64, error) {
	if err := params.Validate(); err != nil {
		return 0, err
	}

	const query = `
    INSERT INTO account_credentials (
        credential,
        account_id,
        provider_id,
        secret,
        secret_type,
        verified,
        verified_at
    )
    VALUES (
        @credential,
        @account_id,
        @provider_id,
        @secret,
        @secret_type,
        @verified,
        @verified_at
    )
    RETURNING id
  `

	var id int64
	err := db.SelectOne(
		ctx,
		query,
		pgx.StrictNamedArgs{
			"credential":  params.Credential,
			"account_id":  params.AccountID,
			"provider_id": params.ProviderID,
			"secret":      params.Secret,
			"secret_type": params.SecretType,
			"verified":    params.Verified,
			"verified_at": params.VerifiedAt,
		},
	).Scan(&id)
	if err != nil {
		return 0, errs.DBErrToAppErr(db.MapDBErrToDomainErr(err))
	}

	return id, nil
}

func (r *CredentialRepository) GetCredentialByProvider(
	ctx context.Context,
	db database.DBTX,
	providerID int64,
	credential string,
) (*model.AccountCredential, error) {
	const query = `
	SELECT 
			id, 
			credential,
			account_id, 
			provider_id, 
			secret,
			secret_type,
			verified,
			verified_at,
			disabled,
			disabled_at,
			created,
			updated 
		FROM account_credentials
	WHERE provider_id = $1 AND credential = $2  
 	`

	var m model.AccountCredential
	if err := db.SelectOne(ctx, query, providerID, credential).Scan(
		&m.ID,
		&m.Credential,
		&m.AccountID,
		&m.ProviderID,
		&m.Secret,
		&m.SecretType,
		&m.Verified,
		&m.VerifiedAt,
		&m.Disabled,
		&m.DisabledAt,
		&m.Created,
		&m.Updated,
	); err != nil {
		return nil, errs.DBErrToAppErr(db.MapDBErrToDomainErr(err))
	}

	return &m, nil
}

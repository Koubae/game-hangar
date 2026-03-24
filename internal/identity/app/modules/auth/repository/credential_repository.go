package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/pkg/database"
)

type ICredentialRepository interface {
	GetCredentialByProvider(ctx context.Context, db database.DBTX, providerID int64, credential string) (*model.AccountCredential, error)
}

// TODO: GetCredentialBy Provider + credential (string) ✅
// TODO: (Service layer..redundant on repo) Does Credential Exists? (By Provider + credential)
// TODO: Get Credential by account_id, provider_id => credential (string)
// TODO: Create Credential
// TODO: Handle Auth
// TODO: CreateCredential
type CredentialRepository struct{}

func NewCredentialRepository() *CredentialRepository {
	r := &CredentialRepository{}
	return r
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("error while GetCredentialByProvider, error: %w", err)
	}

	return &m, nil
}

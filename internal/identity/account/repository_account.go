package account

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/pkg/database"
)

type IAccountRepository interface {
	CreateAccount(
		ctx context.Context,
		db database.DBTX,
		params NewAccount,
	) (*string, error)
	GetAccount(
		ctx context.Context,
		db database.DBTX,
		id string,
	) (*Account, error)
}

type AccountRepositoryFactory func() IAccountRepository

type AccountRepository struct{}

func NewAccountRepository() IAccountRepository {
	r := &AccountRepository{}
	return r
}

func (r *AccountRepository) CreateAccount(
	ctx context.Context,
	db database.DBTX,
	params NewAccount,
) (*string, error) {
	const query = `
		INSERT into account (
				id, 
				username, 
				email
		)
		VALUES (
			gen_random_uuid(),
			@username,
			@email 	
		)
		RETURNING id::text
	`

	var id string
	err := db.SelectOne(
		// TODO: rename this. maybe we should keep same naming convention as pgx API'???
		ctx,
		query,
		pgx.StrictNamedArgs{
			"username": params.Username,
			"email":    params.Email,
		},
	).
		Scan(&id)
	if err != nil {
		return nil, errs.DBErrToAppErr(db.MapDBErrToDomainErr(err), "account")
	}

	return &id, nil
}

func (r *AccountRepository) GetAccount(
	ctx context.Context,
	db database.DBTX,
	id string,
) (*Account, error) {
	const query = `
	SELECT 
			id::text,
			username, 
			email, 
			disabled,
			created, 
			updated 
		FROM account 
	WHERE id = @id 
	`

	var m Account
	if err := db.SelectOne(ctx, query, pgx.StrictNamedArgs{"id": id}).Scan(
		&m.ID,
		&m.Username,
		&m.Email,
		&m.Disabled,
		&m.Created,
		&m.Updated,
	); err != nil {
		return nil, errs.DBErrToAppErr(db.MapDBErrToDomainErr(err), "account")
	}

	return &m, nil
}

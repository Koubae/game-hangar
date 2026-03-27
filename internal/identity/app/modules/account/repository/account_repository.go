package repository

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/model"
	"github.com/koubae/game-hangar/pkg/database"
)

var (
	ErrUsernameRequired   = errors.New("username is required")
	ErrInvalidEmailFormat = errors.New("invalid email format or list")
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
	) (*model.Account, error)
}

type AccountRepositoryFactory func() IAccountRepository

// TODO: This stuff should go in a "domain" layer. or dto??
type NewAccount struct {
	Username string
	Email    *string
}

func (p *NewAccount) Validate() error {
	if strings.TrimSpace(p.Username) == "" {
		return ErrUsernameRequired
	}

	if p.Email != nil {
		_, err := mail.ParseAddressList(*p.Email)
		if err != nil {
			return ErrInvalidEmailFormat
		}
	}

	return nil
}

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
	if err := params.Validate(); err != nil {
		return nil, err
	}

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
	err := db.SelectOne( // TODO: rename this. maybe we should keep same naming convention as pgx API'???
		ctx,
		query,
		pgx.StrictNamedArgs{
			"username": params.Username,
			"email":    params.Email,
		},
	).
		Scan(&id)
	if err != nil {
		return nil, db.MapDBErrToDomainErr(err)
	}

	return &id, nil
}

func (r *AccountRepository) GetAccount(
	ctx context.Context,
	db database.DBTX,
	id string,
) (*model.Account, error) {
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

	var m model.Account
	if err := db.SelectOne(ctx, query, pgx.StrictNamedArgs{"id": id}).Scan(
		&m.ID,
		&m.Username,
		&m.Email,
		&m.Disabled,
		&m.Created,
		&m.Updated,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf( // TODO: Is there a way to automatically "inject" the name ?? of file? or even better struct owner?
			"error while GetAccount, error: %w",
			err,
		)
	}

	return &m, nil
}

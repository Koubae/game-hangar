package account

import (
	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/pkg/web"
)

type DTOAccount struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	Email    *string      `json:"email"`
	Disabled bool         `json:"disabled"`
	Created  web.JSONTime `json:"created"`
	Updated  web.JSONTime `json:"updated"`
}

func NewDTOAccountFromAccount(_account *Account) *DTOAccount {
	return &DTOAccount{
		ID:       _account.ID,
		Username: _account.Username,
		Email:    _account.Email,
		Disabled: _account.Disabled,
		Created:  web.JSONTime(_account.Created),
		Updated:  web.JSONTime(_account.Updated),
	}
}

type DTOCreateAccount struct {
	Source   string `json:"source"   binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (dto *DTOCreateAccount) Validate() *errs.AppError {
	err := errs.DTOSchemaValidation(dto)
	if err != nil {
		return err
	}
	return nil
}

type DTOLoginByUsername struct {
	Source   string `json:"source"   binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (dto *DTOLoginByUsername) Validate() *errs.AppError {
	err := errs.DTOSchemaValidation(dto)
	if err != nil {
		return err
	}
	return nil
}

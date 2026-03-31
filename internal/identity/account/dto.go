package account

import (
	"time"

	"github.com/koubae/game-hangar/internal/errs"
)

type DTOAccount struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

func NewDTOAccountFromAccount(_account *Account) *DTOAccount {
	return &DTOAccount{
		ID:       _account.ID,
		Username: _account.Username,
		Email:    *_account.Email,
		Disabled: _account.Disabled,
		Created:  _account.Created,
		Updated:  _account.Updated,
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

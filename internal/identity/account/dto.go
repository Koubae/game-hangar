package account

import (
	"time"

	"github.com/koubae/game-hangar/pkg/errspkg"
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
	var email string
	if _account.Email != nil {
		email = *_account.Email
	}

	return &DTOAccount{
		ID:       _account.ID,
		Username: _account.Username,
		Email:    email,
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

func (dto DTOCreateAccount) Validate() error {
	err := errspkg.DTOSchemaValidation(dto)
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

func (dto DTOLoginByUsername) Validate() error {
	err := errspkg.DTOSchemaValidation(dto)
	if err != nil {
		return err
	}
	return nil
}

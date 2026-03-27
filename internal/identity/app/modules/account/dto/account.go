package dto

import (
	"errors"
	"fmt"

	"github.com/koubae/game-hangar/pkg/common"
)

type CreateAccountDTO struct {
	Source   string `json:"source"   binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (dto *CreateAccountDTO) Validate() error {
	if err := common.DTOSchemaValidation(dto); err != nil {
		return errors.New(fmt.Sprintf("invalid payload: %v", err))
	}
	return nil
}

type DTOAccount struct {
	AccountID string `json:"account_id"    binding:"required"`
	CredID    int64  `json:"credential_id" binding:"required"`
	Username  string `json:"username"      binding:"required"`
}

package account

import (
	"github.com/koubae/game-hangar/internal/errs"
)

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

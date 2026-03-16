package service

import (
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/app/modules/account/dto"
	"github.com/koubae/game-hangar/pkg/common"
)

type AccountService struct{}

func (s *AccountService) CreateAccount(dto dto.CreateAccountDTO) error {
	err := dto.Validate()
	if err != nil {
		return &common.BusinessError{
			HTTPCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	return nil
}

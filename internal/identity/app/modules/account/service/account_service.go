package service

import (
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/app/modules/account/dto"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/model"
	"github.com/koubae/game-hangar/pkg/common"
)

type AccountService struct{}

func (s *AccountService) CreateAccount(dto dto.CreateAccountDTO) (*model.Account, error) {
	err := dto.Validate()
	if err != nil {
		return nil, &common.BusinessError{
			HTTPCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	// Create Provider Service
	// 		IsProviderEnabled?
	// Create Credential Service
	// 		CredentialExists? 
	
  // 1. Check provider exists 
  // 2. Check whether Acount exists 
  // 3. Check credential exists  
  // 4. .. start tx (transaction)
  // 5. .. create account 
	// 6. create credentials 

	account := model.Account{
		ID:       "uuid-temp",
		Username: dto.Username,
	}
	return &account, nil
}

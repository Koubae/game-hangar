package handlers

import (
	"github.com/koubae/game-hangar/account/internal/domain/account/model"
)

type GetAccountRequest struct {
	FullProfile bool `form:"full_profile" json:"full_profile"`
	Username    string
	ClientID    string
	UserID      uint
}

type GetAccountResponse struct {
	model.Account
	ClientID string
}

type GetAccountHandler struct {
	Command  GetAccountRequest
	Response GetAccountResponse
}

func (h *GetAccountHandler) Handle() error {
	// TODO: Add you know.. any kind of logic here.. database maybe?

	account := model.Account{
		UserID:   h.Command.UserID,
		Username: h.Command.Username,
	}

	h.Response = GetAccountResponse{
		Account:  account,
		ClientID: h.Command.ClientID,
	}
	return nil
}

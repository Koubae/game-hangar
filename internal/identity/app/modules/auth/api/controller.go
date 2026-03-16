package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/app/modules/account/dto"
	accountService "github.com/koubae/game-hangar/internal/identity/app/modules/account/service"
	"github.com/koubae/game-hangar/pkg/common"
)

type AuthController struct{}

func (controller *AuthController) RegisterByUsername(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateAccountDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("invalid json: %v", err), http.StatusBadRequest)
		return
	}

	service := accountService.AccountService{}
	err := service.CreateAccount(payload)

	if err != nil {
		if businessError, ok := errors.AsType[*common.BusinessError](err); ok {
			// TODO: implement API errors
			http.Error(w, businessError.Message, businessError.HTTPCode)
			return
		}
		// TODO: implement API errors
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok RegisterByUsername\n")
}

func (controller *AuthController) LoginByUsername(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok LoginByUsername\n")
}

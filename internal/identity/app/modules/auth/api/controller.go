package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/app/modules/account/dto"
	accountService "github.com/koubae/game-hangar/internal/identity/app/modules/account/service"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/web"
)

type AuthController struct{}

func (controller *AuthController) RegisterByUsername(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateAccountDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.WriteBusinessErrorResponse(
			w, &common.BusinessError{
				HTTPCode: http.StatusBadRequest,
				Message:  fmt.Sprintf("invalid json: %v", err),
			},
		)
		return
	}

	service := accountService.AccountService{}
	account, err := service.CreateAccount(payload)
	if err != nil {
		web.WriteBusinessErrorResponse(w, err)
		return
	}

	web.WriteJSONResponse(w, http.StatusCreated, account)

}

func (controller *AuthController) LoginByUsername(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok LoginByUsername\n")
}

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/app/container"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/dto"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/web"
	"go.uber.org/zap"
)

type AuthController struct {
	container container.IdentityContainer
}

func NewAuthController(container container.IdentityContainer) *AuthController {
	return &AuthController{
		container: container,
	}
}

func (c *AuthController) RegisterByUsername(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	ctx := r.Context()
	logger := c.container.Logger()
	logger.Info(
		"RegisterByUsername called",
		zap.String("source", payload.Source),
		zap.String("username", payload.Username),
	)

	// TODO: Remove this ---------

	secret := payload.Password // TODO: HASHHHHHHHHH
	accountAuthSrv := c.container.AccountAuthService(nil)

	accountID, credID, err := accountAuthSrv.RegisterByUsername(
		ctx,
		payload.Source,
		payload.Username,
		secret,
	)
	if err != nil {
		logger.Error("error while registring account", zap.Error(err))
	}

	// TODO: -----------------------------

	response := dto.DTOAccount{
		AccountID: *accountID,
		CredID:    *credID,
		Username:  payload.Username,
	}
	web.WriteJSONResponse(w, http.StatusCreated, response)
}

func (c *AuthController) LoginByUsername(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok LoginByUsername\n")
}

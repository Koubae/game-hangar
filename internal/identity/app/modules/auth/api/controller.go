package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account"
	"github.com/koubae/game-hangar/internal/identity/container"
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
	var payload account.CreateAccountDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.WriteBusinessErrorResponse(
			w, &common.ClientResponseError{
				HTTPCode: http.StatusBadRequest,
				Message:  fmt.Sprintf("invalid json: %v", err.Error()),
			},
		)
		return
	}

	if err := payload.Validate(); err != nil {
		web.WriteBusinessErrorResponse(
			w, &common.ClientResponseError{
				HTTPCode: http.StatusBadRequest,
				Message:  fmt.Sprintf("invalid payload: %s", err.Error()),
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

	secret, err := c.container.AuthService().HashSecret(payload.Password)
	if err != nil {
		response := errs.AppErrToClientResponseWithLog(
			err,
			"error while hashing secret during registration by username",
			logger,
		)
		web.WriteBusinessErrorResponse(w, response)
		return
	}

	accountAuthSrv := c.container.AccountAuthService(nil)
	accountID, credID, err := accountAuthSrv.RegisterByUsername(
		ctx,
		payload.Source,
		payload.Username,
		secret,
	)
	if err != nil {
		response := errs.AppErrToClientResponseWithLog(err, "could not create account", logger)
		web.WriteBusinessErrorResponse(w, response)
		return
	}

	response := account.DTOAccount{
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

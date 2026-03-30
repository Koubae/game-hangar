package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/internal/errs"
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
				Message:  fmt.Sprintf("invalid json: %v", err.Error()),
			},
		)
		return
	}

	if err := payload.Validate(); err != nil {
		web.WriteBusinessErrorResponse(
			w, &common.BusinessError{
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
		logger.Error(
			"error while hashing secret during registration by username",
			zap.Error(err),
		)
		web.WriteBusinessErrorResponse(
			w, &common.BusinessError{
				HTTPCode: http.StatusInternalServerError,
				Message:  "unexpected error occurred",
			},
		)
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
		var responseError *common.BusinessError
		var lvl string
		appErr := errs.AsAppError(err)
		if appErr.IsServerErr() {
			lvl = "error"
			responseError = &common.BusinessError{
				HTTPCode: http.StatusInternalServerError,
				Message:  "unexpected error occurred",
			}
		} else {
			lvl = "info"
			responseError = &common.BusinessError{
				HTTPCode: http.StatusBadRequest,
				Message: fmt.Sprintf(
					"could not create account, error: %s",
					err.Error(),
				),
			}
		}

		logger.L(lvl, "could not create account", zap.Error(err))
		web.WriteBusinessErrorResponse(w, responseError)
		return
	}

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

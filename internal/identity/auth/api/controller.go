package api

import (
	"io"
	"net/http"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/internal/identity/container"
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
	payload, ok := web.LoadJsonBodyOrBadRequestResponse[account.DTOCreateAccount](w, r)
	if !ok {
		return
	}
	if err := payload.Validate(); err != nil {
		errs.AppErrToClientResponse(w, err, "")
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
		errs.AppErrToClientResponseWithLog(
			w,
			err,
			"hash secret error on registration by username",
			logger,
		)
		return
	}

	accountID, credID, err := c.container.AccountAuthService(nil).RegisterByUsername(
		ctx,
		payload.Source,
		payload.Username,
		secret,
	)
	if err != nil {
		errs.AppErrToClientResponseWithLog(w, err, "could not create account", logger)
		return
	}

	response := auth.DTOAccountLoggedIn{
		AccountID:    *accountID,
		Username:     payload.Username,
		LoggedCredID: *credID,
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

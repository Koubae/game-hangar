package api

import (
	"net/http"
	"time"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/web"
	"go.uber.org/zap"
)

const (
	AuthTokenExpirationTime = time.Hour * 4
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
	payload, ok := web.LoadJsonBody[account.DTOCreateAccount](w, r)
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

	secretService := c.container.SecretsService()
	err := secretService.ValidatePasswordDefaultRules(payload.Password)
	if err != nil {
		errs.AppErrToClientResponseWithLog(w, err, "", logger)
		return
	}

	secret, err := secretService.HashSecret(payload.Password)
	if err != nil {
		errs.AppErrToClientResponseWithLog(w, err, "hash secret error on registration by username", logger)
		return
	}

	accountID, credID, err := c.container.AccountAuthService(nil).RegisterByUsername(
		ctx,
		payload.Source,
		payload.Username,
		secret,
	)
	if err != nil {
		errs.AppErrToClientResponseWithLog(w, err, "could not create account: ", logger)
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
	payload, ok := web.LoadJsonBody[account.DTOLoginByUsername](w, r)
	if !ok {
		return
	}
	if err := payload.Validate(); err != nil {
		errs.AppErrToClientResponse(w, err, "")
		return
	}

	ctx := r.Context()
	logger := c.container.Logger()
	logger.Debug(
		"LoginByUsername called",
		zap.String("source", payload.Source),
		zap.String("username", string(auth.Username)),
	)

	provider, err := c.container.ProviderService(nil).GetEnabledProvider(ctx, payload.Source, string(auth.Username))
	if err != nil {
		errs.AppErrToClientResponseWithLog(w, errs.Wrap(errs.AuthLoginFailed, err), "", logger)
		return
	}

	credential, err := c.container.CredentialService(nil).GetCredentialByProvider(ctx, provider.ID, payload.Username)
	if err != nil {
		errs.AppErrToClientResponseWithLog(w, errs.Wrap(errs.AuthLoginFailed, err), "credential ", logger)
		return
	}

	secretService := c.container.SecretsService()
	if !secretService.VerifySecret(credential.Secret, payload.Password) {
		errs.AppErrToClientResponseWithLog(
			w,
			errs.Wrap(errs.AuthLoginFailed, errs.AuthLoginPasswordMismatch),
			"credential ",
			logger,
		)
		return
	}

	expire := time.Now().Add(AuthTokenExpirationTime).Unix()
	accessToken, err := secretService.GenerateJWTAccessToken(
		provider.Source,
		provider.Type,
		credential.AccountID.String(),
		credential.Credential,
		expire,
	)

	response := auth.DTOAccessToken{
		AccessToken: accessToken,
		ExpiresIn:   expire,
	}
	web.WriteJSONResponse(w, http.StatusOK, response)

}

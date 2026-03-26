package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/app/container"
	"github.com/koubae/game-hangar/internal/identity/app/modules/account/dto"
	accountService "github.com/koubae/game-hangar/internal/identity/app/modules/account/service"
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
		zap.String("username", payload.Username),
	)

	// TODO: Remove this ---------
	providerRepo := c.container.ProviderRepository()
	logger.Info(
		"provider repo",
		zap.String("repo", fmt.Sprintf("%+v", providerRepo)),
	)

	credRepo := c.container.CredentialRepository()
	logger.Info(
		"cre repo",
		zap.String("credRepo", fmt.Sprintf("%v", credRepo)),
	)

	providerID := 1
	credential := "account_test_1"
	cred, err := credRepo.GetCredentialByProvider(
		ctx,
		c.container.DB(),
		int64(providerID),
		credential,
	)
	if err != nil {
		logger.Warn(
			"error while gett cred",
			zap.String("cred", credential),
			zap.Error(err),
		)
	} else {
		logger.Info("succcess cred",
			zap.String("cred", cred.Credential), zap.String("accID", cred.AccountID.String()))
	}

	providerService := c.container.ProviderService(nil)
	isUsernameAuthEnabled := providerService.IsProviderEnabled(
		ctx,
		"global",
		"username",
	)
	logger.Info(
		"is username prov enabled?",
		zap.Bool("enabled?", isUsernameAuthEnabled),
	)

	// TODO: -----------------------------

	service := accountService.AccountService{}
	account, err := service.CreateAccount(payload)
	if err != nil {
		web.WriteBusinessErrorResponse(w, err)
		return
	}

	response := dto.DTOAccount{
		ID:       account.ID,
		Username: account.Username,
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

package api

import (
	"fmt"
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"github.com/koubae/game-hangar/pkg/web"
	"go.uber.org/zap"
)

type AccountManagementController struct {
	container container.IdentityContainer
}

func NewAccountManagementController(c container.IdentityContainer) *AccountManagementController {
	return &AccountManagementController{
		container: c,
	}
}

func (c *AccountManagementController) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken, ok := authpkg.GetAccessToken(ctx)
	if !ok {
		errspkg.AppErrToClientResponse(w, errspkg.AuthNotLoggedIn, "")
		return
	}

	logger := c.container.Logger()

	// TODO: rem -- dev
	permissions := authpkg.GetPermissionsOrDefault(ctx)
	logger.Info("permissions", zap.String("permissions", fmt.Sprintf("%v", permissions)))
	logger.Info("access_token issuer", zap.String("issuer", accessToken.Issuer))
	// TODO: rem -- dev

	_account, err := c.container.AccountManagementService(nil).GetAccount(
		ctx,
		accessToken.AccountID,
	)
	if err != nil {
		errspkg.AppErrToClientResponseWithLog(w, err, "", logger)
		return
	}

	response := account.NewDTOAccountFromAccount(_account)
	web.WriteJSONResponse(w, http.StatusOK, response)
}

package api

import (
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/di"
	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(v1 *http.ServeMux, c di.Container) {
	loggedAccountMiddleware := authpkg.NewJWTMiddleware()

	account := web.Group(v1, "/account")
	accountLoggedIn := web.GroupWithMiddleware(
		account,
		"",
		loggedAccountMiddleware,
	)

	authController := NewAccountManagementController(c.(container.IdentityContainer))

	// ------------------------------------------
	// 	Account Management functions | Account Access
	// ------------------------------------------
	accountLoggedIn.HandleFunc("GET /me", authController.Me)

	// ------------------------------------------
	// 	Account Management functions | Backoffice Access
	// ------------------------------------------
	loggedAdminMiddleware := authpkg.NewAdminJWTMiddleware()
	accountBackofficeLoggedIn := web.GroupWithMiddleware(
		v1,
		"/backoffice/account",
		loggedAdminMiddleware,
	)
	accountBackofficeLoggedIn.HandleFunc("GET /{id}", authController.GetAccount)
}

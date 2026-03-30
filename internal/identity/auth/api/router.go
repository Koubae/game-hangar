package api

import (
	"net/http"

	"github.com/koubae/game-hangar/internal/identity/container"
	"github.com/koubae/game-hangar/pkg/di"
	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(v1 *http.ServeMux, c di.Container) {
	auth := web.Group(v1, "/auth")

	authController := NewAuthController(c.(container.IdentityContainer))

	// ------------------------------------------
	// 	Register functions
	// ------------------------------------------
	auth.HandleFunc(
		"POST /register/username",
		authController.RegisterByUsername,
	)

	// ------------------------------------------
	// 	Login functions
	// ------------------------------------------
	auth.HandleFunc("POST /login/username", authController.LoginByUsername)
}

package api

import (
	"net/http"

	"github.com/koubae/game-hangar/pkg/di"
	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(v1 *http.ServeMux, container *di.Container) {
	auth := web.Group(v1, "/auth")

	authController := NewAuthController(container)

	// ------------------------------------------
	// 	Register functions
	// ------------------------------------------
	auth.HandleFunc("POST /register/username", authController.RegisterByUsername)

	// ------------------------------------------
	// 	Login functions
	// ------------------------------------------
	auth.HandleFunc("POST /login/username", authController.LoginByUsername)

}

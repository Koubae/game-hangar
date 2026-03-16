package api

import (
	"net/http"

	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(v1 *http.ServeMux) {
	auth := web.Group(v1, "/auth")

	authController := new(AuthController)

	// ------------------------------------------
	// 	Register functions
	// ------------------------------------------
	auth.HandleFunc("POST /register/username", authController.RegisterByUsername)

	// ------------------------------------------
	// 	Login functions
	// ------------------------------------------
	auth.HandleFunc("POST /login/username", authController.LoginByUsername)

}

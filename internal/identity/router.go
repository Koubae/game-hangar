package identity

import (
	"fmt"
	"io"
	"net/http"

	authRouter "github.com/koubae/game-hangar/internal/identity/auth/api"
	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/di"
	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(container di.Container) web.RouterRegisterFunc {
	return func(mux *http.ServeMux) {
		loggedAccountMiddleware := authpkg.NewJWTMiddleware()

		v1 := web.Group(mux, "/api/v1")

		authRouter.RouterRegister(v1, container)

		account := web.Group(v1, "/account")
		account.HandleFunc(
			"GET /temp", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := io.WriteString(w, fmt.Sprintf("cool beans"))
				if err != nil {
					return
				}
			},
		)

		protected := web.GroupWithMiddleware(
			account,
			"/protected",
			loggedAccountMiddleware,
		)

		protected.HandleFunc(
			"/me", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("protected content"))
			},
		)

		// account.Handle(
		// 	"/protected",
		// 	auth.JWTMiddleware(
		// 		jwt.SigningMethodHS256, secret, http.HandlerFunc(
		// 			func(w http.ResponseWriter, r *http.Request) {
		// 				w.Write([]byte("ok"))
		// 			},
		// 		),
		// 	),
		// )

	}

}

package identity

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	authRouter "github.com/koubae/game-hangar/internal/identity/auth/api"
	"github.com/koubae/game-hangar/pkg/auth"
	"github.com/koubae/game-hangar/pkg/di"
	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(container di.Container) web.RouterRegisterFunc {
	return func(mux *http.ServeMux) {

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

		// var secret []byte = []byte("your-secret")
		secret, _ := jwt.ParseRSAPublicKeyFromPEM([]byte("your-secret"))
		protected := web.GroupWithMiddleware(
			account,
			"/protected",
			auth.JWTMiddleware[*rsa.PublicKey](jwt.SigningMethodHS256, secret),
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

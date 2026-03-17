package app

import (
	"fmt"
	"io"
	"net/http"

	authRouter "github.com/koubae/game-hangar/internal/identity/app/modules/auth/api"
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

	}

}

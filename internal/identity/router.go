package identity

import (
	"net/http"

	accountRouter "github.com/koubae/game-hangar/internal/identity/account/api"
	authRouter "github.com/koubae/game-hangar/internal/identity/auth/api"
	"github.com/koubae/game-hangar/pkg/di"
	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(container di.Container) web.RouterRegisterFunc {
	return func(mux *http.ServeMux) {
		v1 := web.Group(mux, "/api/v1")

		authRouter.RouterRegister(v1, container)
		accountRouter.RouterRegister(v1, container)
	}

}

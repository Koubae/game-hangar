package api

import (
	"io"
	"net/http"

	"github.com/koubae/game-hangar/pkg/web"
)

func RouterRegister(v1 *http.ServeMux) {
	auth := web.Group(v1, "/auth")
	auth.HandleFunc(
		"POST /authorize", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, "ok 2")
		},
	)
	auth.HandleFunc(
		"POST /login", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, "ok")
		},
	)

}

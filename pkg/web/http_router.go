package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/middleware"
	"github.com/rs/cors"
)

type RouterRegisterFunc func(mux *http.ServeMux)
type RouterFunc func(logger common.Logger, config *common.Config, routerRegister RouterRegisterFunc) *http.Handler

func Router(logger common.Logger, config *common.Config, routerRegister RouterRegisterFunc) *http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"GET /{$}", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, fmt.Sprintf("Welcome to %s", config.GetFullName()))
			if err != nil {
				return
			}
		},
	)

	routerRegister(mux)

	handler := cors.New(
		cors.Options{
			AllowedOrigins:   config.CORSConfig.AllowedOrigins,
			AllowedMethods:   config.CORSConfig.AllowedMethods,
			AllowedHeaders:   config.CORSConfig.AllowedHeaders,
			AllowCredentials: config.CORSConfig.AllowCredentials,
		},
	).Handler(mux)
	handler = middleware.AccessLogger(logger, handler)
	handler = middleware.RecoveryMiddleware(logger, handler)
	return &handler

}

func Group(mux *http.ServeMux, prefix string) *http.ServeMux {
	sub := http.NewServeMux()
	mux.Handle(prefix+"/", http.StripPrefix(prefix, sub))
	return sub
}

func GroupWithMiddleware(mux *http.ServeMux, prefix string, mw func(http.Handler) http.Handler) *http.ServeMux {
	sub := http.NewServeMux()
	mux.Handle(prefix+"/", mw(http.StripPrefix(prefix, sub)))
	return sub
}

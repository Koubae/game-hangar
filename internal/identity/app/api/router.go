package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/middleware"
	"github.com/rs/cors"
)

func Router(logger common.Logger, config *common.Config) *http.Handler {
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

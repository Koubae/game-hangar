package app

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

func RunServer() {
	// TODO: load .env file
	logLevel := "info"

	logger := common.CreateLogger(logLevel)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	logger.Info(
		"server started",
		zap.String("addr", ":8080"),
		zap.String("env", "dev"),
	)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /{$}", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, "Hello World :)\n")
			if err != nil {
				return
			}
		},
	)

	srv := http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		// Use CrossOriginProtection.Handler to block all non-safe cross-origin
		// browser requests to mux.
		Handler: http.NewCrossOriginProtection().Handler(mux),
	}

	log.Fatal(srv.ListenAndServe())
}

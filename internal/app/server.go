package app

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

func RunServer() {
	// TODO: load .env file
	logLevel := "info"
	// -----------

	logger := common.CreateLogger(logLevel)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

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

	addr := ":8080"

	srv := http.Server{
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		// Use CrossOriginProtection.Handler to block all non-safe cross-origin
		// browser requests to mux.
		Handler: http.NewCrossOriginProtection().Handler(mux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while shutting down the server, error: %s", err)
		}
	}()

	logger.Info(
		"server started",
		zap.String("addr", addr),
		zap.String("env", "dev"),
	)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	logger.Info("Shutdown signal received, shutting down the server")

	// clean up

	logger.Info("Server shutdown completed")
}

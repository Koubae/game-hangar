package app

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

const SHUTDOWN_GRACEFULLY_TIMEOUT_SECONDS = 10 * time.Second

func RunServer() {
	loggerTmp, _ := common.CreateLogger(common.LogLevelInfo, "")
	config := NewConfig(loggerTmp)

	logger, logCloser := common.CreateLogger(config.LogLevel, config.LogFilePath)
	defer logCloser(loggerTmp, logger)

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

	handler := cors.New(
		cors.Options{
			AllowedOrigins:   config.CORSConfig.AllowedOrigins,
			AllowedMethods:   config.CORSConfig.AllowedMethods,
			AllowedHeaders:   config.CORSConfig.AllowedHeaders,
			AllowCredentials: config.CORSConfig.AllowCredentials,
		},
	).Handler(mux)

	srv := http.Server{
		Addr:           config.GetAppURL(),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 5120,
		Handler:        handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while shutting down the server, error: %s", err)
		}
	}()

	logger.Info(
		"Server started",
		zap.String("addr", config.GetAppURL()),
		zap.String("env", string(config.Env)),
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP,
	)
	defer stop()

	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify the user of shutdown.
	stop()

	logger.Info(
		"Shutdown signal received, Shutting down server gracefully... ",
		zap.Duration("shutdown_gracefull_timeout", SHUTDOWN_GRACEFULLY_TIMEOUT_SECONDS),
	)
	ctx, cancel := context.WithTimeout(context.Background(), SHUTDOWN_GRACEFULLY_TIMEOUT_SECONDS)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown Failed", zap.Error(err))
	}

	logger.Info("Server has shutdown, cleaning up resources ...")

	// TODO: clean up resources here...

	logger.Info("Resource cleanup completed, terminating process...")
}

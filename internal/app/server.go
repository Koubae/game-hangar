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
		ReadTimeout:    time.Duration(config.ServerReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.ServerWriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.ServerIdleTimeout) * time.Second,
		MaxHeaderBytes: config.ServerMaxHeaderBytes,
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

	serverShutdownGraceTimeout := time.Duration(config.ServerShutdownGraceTimeout) * time.Second
	logger.Info(
		"Shutdown signal received, Shutting down server gracefully... ",
		zap.Duration("shutdown_gracefull_timeout", serverShutdownGraceTimeout),
	)
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownGraceTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown Failed", zap.Error(err))
	}

	logger.Info("Server has shutdown, cleaning up resources ...")

	// TODO: clean up resources here...

	logger.Info("Resource cleanup completed, terminating process...")
}

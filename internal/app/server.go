package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koubae/game-hangar/internal/app/api"
	"github.com/koubae/game-hangar/internal/app/settings"
	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

func RunServer() {
	loggerTmp, _ := common.CreateLogger(common.LogLevelInfo, "")
	config := settings.NewConfig(loggerTmp)

	logger, logCloser := common.CreateLogger(config.LogLevel, config.LogFilePath)
	defer logCloser(loggerTmp, logger)

	routerHandler := api.Router(logger)
	srv := http.Server{
		Addr:           config.GetAppURL(),
		ReadTimeout:    time.Duration(config.ServerReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.ServerWriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.ServerIdleTimeout) * time.Second,
		MaxHeaderBytes: config.ServerMaxHeaderBytes,
		Handler:        *routerHandler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while shutting down the server, error: %s", err)
		}
	}()

	logger.Info(
		fmt.Sprintf("Server started on %s -- App: %s", config.GetAppURL(), config.GetFullName()),
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

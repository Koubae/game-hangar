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

type App struct {
	Config    *settings.Config
	Logger    *zap.Logger
	LogCloser func(tmpLogger *zap.Logger, logger *zap.Logger)
	Server    *http.Server
}

func NewApp() *App {
	loggerTmp, _ := common.CreateLogger(common.LogLevelInfo, "")
	config := settings.NewConfig(loggerTmp)

	logger, logCloser := common.CreateLogger(config.LogLevel, config.LogFilePath)

	routerHandler := api.Router(logger, config)
	srv := &http.Server{
		Addr:           config.GetAppURL(),
		ReadTimeout:    time.Duration(config.ServerReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.ServerWriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.ServerIdleTimeout) * time.Second,
		MaxHeaderBytes: config.ServerMaxHeaderBytes,
		Handler:        *routerHandler,
	}

	return &App{
		Config:    config,
		Logger:    logger,
		LogCloser: logCloser,
		Server:    srv,
	}
}

func (a *App) Start() {
	loggerTmp, _ := common.CreateLogger(common.LogLevelInfo, "")
	defer a.LogCloser(loggerTmp, a.Logger)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while shutting down the server, error: %s", err)
		}
	}()

	a.Logger.Info(
		fmt.Sprintf("Server started on %s -- App: %s", a.Config.GetAppURL(), a.Config.GetFullName()),
		zap.String("addr", a.Config.GetAppURL()),
		zap.String("env", string(a.Config.Env)),
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP,
	)
	defer stop()

	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify the user of shutdown.
	stop()

}

func (a *App) Stop() {
	serverShutdownGraceTimeout := time.Duration(a.Config.ServerShutdownGraceTimeout) * time.Second
	a.Logger.Info(
		"Shutdown signal received, Shutting down server gracefully... ",
		zap.Duration("shutdown_gracefull_timeout", serverShutdownGraceTimeout),
	)
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownGraceTimeout)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Fatal("Server Shutdown Failed", zap.Error(err))
	}

	a.Logger.Info("Server has shutdown, cleaning up resources ...")

	// TODO: clean up resources here...

	a.Logger.Info("Resource cleanup completed, terminating process...")
}

func RunServer() {
	app := NewApp()
	app.Start()
	app.Stop()
}

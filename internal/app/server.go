package app

import (
	"context"
	"errors"
	"fmt"
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

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Handler() http.Handler
}

type App struct {
	Config    *settings.Config
	Logger    common.Logger
	LogCloser func(tmpLogger *zap.Logger, logger *zap.Logger)
	Server    Server
}

type httpServerWrapper struct {
	*http.Server
}

func (s *httpServerWrapper) Handler() http.Handler {
	return s.Server.Handler
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
		Server:    &httpServerWrapper{srv},
	}
}

func (a *App) Start(ctx context.Context) {
	loggerTmp, _ := common.CreateLogger(common.LogLevelInfo, "")
	defer func() {
		if z, ok := a.Logger.(*zap.Logger); ok {
			a.LogCloser(loggerTmp, z)
		}
	}()

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	a.Logger.Info(
		fmt.Sprintf("Server started on %s -- App: %s", a.Config.GetAppURL(), a.Config.GetFullName()),
		zap.String("addr", a.Config.GetAppURL()),
		zap.String("env", string(a.Config.Env)),
	)

	sigCtx, stop := signal.NotifyContext(
		ctx,
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP,
	)
	defer stop()

	<-sigCtx.Done()
	// Restore default behavior on the interrupt signal and notify the user of shutdown.
	stop()
}

func (a *App) Stop() error {
	serverShutdownGraceTimeout := time.Duration(a.Config.ServerShutdownGraceTimeout) * time.Second
	a.Logger.Info(
		"Shutdown signal received, Shutting down server gracefully... ",
		zap.Duration("shutdown_gracefull_timeout", serverShutdownGraceTimeout),
	)
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownGraceTimeout)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Error("Server Shutdown Failed", zap.Error(err))
		return err
	}

	a.Logger.Info("Server has shutdown, cleaning up resources ...")

	// TODO: clean up resources here...

	a.Logger.Info("Resource cleanup completed, terminating process...")
	return nil
}

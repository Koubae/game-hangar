package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/di"
	"go.uber.org/zap"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Handler() http.Handler
}

type HTTPApp struct {
	Config    *common.Config
	Server    Server
	Container *di.Container
}

type httpServerWrapper struct {
	*http.Server
}

func (s *httpServerWrapper) Handler() http.Handler {
	return s.Server.Handler
}

func NewHTTPApp(appPrefix string, container *di.Container, config *common.Config, router RouterFunc, routerRegister RouterRegisterFunc) *HTTPApp {
	routerHandler := router(container, config, routerRegister)
	srv := &http.Server{
		Addr:           config.GetAppURL(),
		ReadTimeout:    time.Duration(config.ServerReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.ServerWriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.ServerIdleTimeout) * time.Second,
		MaxHeaderBytes: config.ServerMaxHeaderBytes,
		Handler:        *routerHandler,
	}

	return &HTTPApp{
		Config:    config,
		Server:    &httpServerWrapper{srv},
		Container: container,
	}
}

func (a *HTTPApp) Start(ctx context.Context) {
	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Container.Logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	a.Container.Logger.Info(
		"Server started",
		zap.String("addr", a.Config.GetAppURL()),
		zap.String("app", a.Config.GetFullName()),
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

func (a *HTTPApp) Stop() error {
	defer func() {
		a.Container.Logger.Info("Server has shutdown, cleaning up resources ...")

		if a.Container != nil {
			if err := a.Container.Shutdown(); err != nil {
				a.Container.Logger.Error("Container Shutdown Failed", zap.Error(err))
			}
		}

		a.Container.Logger.Info("Resource cleanup completed, terminating process...")
	}()

	serverShutdownGraceTimeout := time.Duration(a.Config.ServerShutdownGraceTimeout) * time.Second
	a.Container.Logger.Info(
		"Shutdown signal received, Shutting down server gracefully... ",
		zap.Duration("shutdown_gracefull_timeout", serverShutdownGraceTimeout),
	)
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownGraceTimeout)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Container.Logger.Error("Server Shutdown Failed", zap.Error(err))
		return err
	}

	return nil
}

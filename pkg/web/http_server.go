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
	"go.uber.org/zap"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Handler() http.Handler
}

type HTTPApp struct {
	Config *common.Config
	Logger common.Logger
	Server Server
}

type httpServerWrapper struct {
	*http.Server
}

func (s *httpServerWrapper) Handler() http.Handler {
	return s.Server.Handler
}

func NewHTTPApp(appPrefix string, router RouterFunc, routerRegister RouterRegisterFunc) *HTTPApp {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	config := common.NewConfig(loggerTmp, appPrefix)

	logger := common.CreateLogger(config.LogLevel, config.LogFilePath)

	routerHandler := router(logger, config, routerRegister)
	srv := &http.Server{
		Addr:           config.GetAppURL(),
		ReadTimeout:    time.Duration(config.ServerReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.ServerWriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.ServerIdleTimeout) * time.Second,
		MaxHeaderBytes: config.ServerMaxHeaderBytes,
		Handler:        *routerHandler,
	}

	return &HTTPApp{
		Config: config,
		Logger: logger,
		Server: &httpServerWrapper{srv},
	}
}

func (a *HTTPApp) Start(ctx context.Context) {
	loggerTmp := common.CreateLogger(common.LogLevelInfo, "")
	defer func() {
		if z, ok := a.Logger.(*common.AppLogger); ok {
			z.LogCloser(loggerTmp, z.Logger)
		}
	}()

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	a.Logger.Info(
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

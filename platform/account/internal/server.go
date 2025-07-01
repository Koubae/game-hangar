package accountpackage

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/koubae/game-hangar/account/internal/infrastructure/api/routes"
	"github.com/koubae/game-hangar/account/internal/settings"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	config := settings.NewConfig()
	switch config.Environment {
	case settings.EnvTesting:
		gin.SetMode(gin.TestMode)
	case settings.EnvDev, settings.EnvStaging:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}

func RunServer() {
	config := settings.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	err := router.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		panic(err.Error())
	}

	routes.InitRoutes(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while shutting down server, error: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify the user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

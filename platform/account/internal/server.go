package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/koubae/game-hangar/account/internal/infrastructure/api/routes"
	"github.com/koubae/game-hangar/account/internal/settings"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	settings.NewConfig()
}

func RunServer() {
	config := settings.GetConfig()

	switch config.Environment {
	case settings.EnvTesting:
		gin.SetMode(gin.TestMode)
	case settings.EnvDev, settings.EnvStaging:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	err := router.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		panic(err.Error())
	}

	routes.InitRoutes(router)

	// TODO: Graceful shutdown
	err = router.Run()
	if err != nil {
		panic(err.Error())
	}
}

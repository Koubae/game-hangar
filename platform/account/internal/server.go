package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/koubae/game-hangar/account/internal/infrastructure/api/routes"
	"os"
)

func RunServer() {
	errTemp := os.Setenv("PORT", "8001") // TODO: configurable, is Go/gin var
	if errTemp != nil {
		panic(errTemp.Error())
	}
	gin.SetMode(gin.DebugMode) // TODO: Configurable

	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.2"}) // TODO: just an example!
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

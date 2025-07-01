package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koubae/game-hangar/account/internal/infrastructure/api/controllers"
	middlewares2 "github.com/koubae/game-hangar/account/internal/infrastructure/api/middlewares"
	"github.com/koubae/game-hangar/account/pkg/utils"
)

// TODO: .env variable!

func InitRoutes(router *gin.Engine) {
	index := router.Group("/")
	{
		index.GET("/", func(c *gin.Context) {
			c.String(200, "Service Running...")
		})

		index.GET("/health", func(c *gin.Context) {
			c.String(200, "OK")
		})

		index.GET("/ready", func(c *gin.Context) {
			c.String(200, "OK")
		})
	}

	v1 := router.Group("/api/v1")

	authController := controllers.AuthController{}
	authV1 := v1.Group("/auth")
	{
		authV1.POST("/login", authController.LoginV1)
	}

	accountControllers := controllers.AccountControllers{}
	// TODO - load on config on stasrt up!
	publicKey := utils.GetPublicKeyOrPanic()
	accountV1 := v1.Group("/account", middlewares2.IsAuthorizedJWTWithRSAMiddleware(publicKey))
	{
		accountV1.POST("", accountControllers.Create)
		accountV1.GET("/:name", accountControllers.Get)
	}
}

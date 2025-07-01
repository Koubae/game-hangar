package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koubae/game-hangar/account/internal/infrastructure/api/controllers"
	middlewares2 "github.com/koubae/game-hangar/account/internal/infrastructure/api/middlewares"
)

func InitRoutes(router *gin.Engine) {
	authMiddleWare := middlewares2.NewJWTRSAMiddleware()

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
	accountV1 := v1.Group("/account", authMiddleWare)
	{
		accountV1.POST("", accountControllers.Create)
		accountV1.GET("/:name", accountControllers.Get)
	}
}

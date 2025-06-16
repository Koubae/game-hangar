package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/koubae/game-hangar/account/internal/infrastructure/api/controllers/account"
)

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

	accountControllers := controllers.AccountControllers{}
	accountV1 := v1.Group("/account")
	{
		accountV1.GET("/:name", accountControllers.Get)
	}

}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koubae/game-hangar/account/internal/domain/api/middlewares"
	"github.com/koubae/game-hangar/account/internal/infrastructure/api/controllers"
	"github.com/koubae/game-hangar/account/pkg/utils"
)

// TODO: .env variable!
var AUTH_SECRET = []byte("AUTH_SECRET_1234")

//var AUTH_SECRET = []byte(os.Getenv("AUTH_JWT_SECRET"))

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
	v2 := router.Group("/api/v2")

	authController := controllers.AuthController{}
	authV1 := v1.Group("/auth")
	{
		authV1.POST("/login", authController.LoginV1)
	}
	authV2 := v2.Group("/auth")
	{
		authV2.POST("/login", authController.LoginV2)
	}

	accountControllers := controllers.AccountControllers{}
	accountV1 := v1.Group("/account", middlewares.IsAuthorizedJWTWithHMACMiddleware(AUTH_SECRET))
	{
		accountV1.POST("", accountControllers.Create)
		accountV1.GET("/:name", accountControllers.Get)
	}

	// TODO - load on config on stasrt up!
	publicKey, err := utils.GetPublicKey()
	if err != nil {
		panic(err.Error())
	}
	// TODO
	accountV2 := v2.Group("/account", middlewares.IsAuthorizedJWTWithRSAMiddleware(publicKey))
	{
		accountV2.POST("", accountControllers.Create)
		accountV2.GET("/:name", accountControllers.Get)
	}
}

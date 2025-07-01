package middlewares

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/koubae/game-hangar/account/pkg/utils"
)

func NewJWTRSAMiddleware() gin.HandlerFunc {
	publicKey := utils.GetPublicKeyOrPanic()
	return func(c *gin.Context) {
		jwtMiddleware[*rsa.PublicKey](c, jwt.SigningMethodRS256, publicKey)
	}
}

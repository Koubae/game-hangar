package middlewares

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/koubae/game-hangar/account/pkg/utils"
)

func IsAuthorizedJWTWithRSAMiddleware(publicKey *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMiddleware(c, jwt.SigningMethodRS256, publicKey)
	}
}

func NewJWTRSAMiddleware() gin.HandlerFunc {
	publicKey := utils.GetPublicKeyOrPanic()
	return IsAuthorizedJWTWithRSAMiddleware(publicKey)
}

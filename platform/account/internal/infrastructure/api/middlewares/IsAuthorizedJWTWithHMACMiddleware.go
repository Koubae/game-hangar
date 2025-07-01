package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthorizedJWTWithHMACMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMiddleware[[]byte](c, jwt.SigningMethodHS256, secret)
	}
}

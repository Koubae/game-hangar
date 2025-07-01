package middlewares

import (
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type JWTSecret interface {
	[]byte | *rsa.PublicKey | *rsa.PrivateKey
}

func jwtMiddleware[S JWTSecret](c *gin.Context, method jwt.SigningMethod, secret S) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != method {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	var exp float64
	exp, ok = claims["exp"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims exp"})
		return
	}

	now := float64(time.Now().Unix())
	if exp < now {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	}

	c.Set("user_id", claims["sub"])
	c.Set("issuer", claims["iss"])
	c.Set("role", claims["role"])
	c.Set("client_id", claims["client_id"])

	c.Next()
}

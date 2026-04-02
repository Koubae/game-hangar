package testutil

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
	"github.com/koubae/game-hangar/pkg/authpkg"
)

func ExtractAdminAccessToken(tokenString string) (*authpkg.AccessToken, error) {
	secret := authpkg.GetAdminPublicKey()
	method := jwt.SigningMethodRS256

	return authpkg.ParseJWToken[*rsa.PublicKey](
		method,
		secret,
		tokenString,
		true,
	)

}

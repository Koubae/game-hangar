package authpkg

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/web"
)

type JWTSecret interface {
	[]byte | *rsa.PublicKey | *rsa.PrivateKey
}

type Middleware func(http.Handler) http.Handler
type contextKey string

const (
	ContextKeySource      contextKey = "source"
	ContextKeyType        contextKey = "type"
	ContextKeyAccountID   contextKey = "account_id"
	ContextKeyCredential  contextKey = "credential"
	ContextKeyIssuer      contextKey = "issuer"
	ContextKeyRole        contextKey = "role"
	ContextKeyAccessToken contextKey = "access_token"
)

func NewJWTMiddleware() func(http.Handler) http.Handler {
	secret := GetPublicKey()
	return JWTMiddleware[*rsa.PublicKey](jwt.SigningMethodRS256, secret)
}

func JWTMiddleware[S JWTSecret](method jwt.SigningMethod, secret S) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				tokenString, ok := extractToken(r)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "missing or invalid token",
						},
					)
					return
				}

				token, err := jwt.Parse(
					tokenString, func(t *jwt.Token) (interface{}, error) {
						if t.Method.Alg() != method.Alg() {
							return nil, fmt.Errorf("unexpected signing method")
						}
						return secret, nil
					},
				)

				if err != nil {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  fmt.Sprintf("token parsing error, %v", err),
						},
					)
					return
				}
				if !token.Valid {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}

				source, ok := claims["source"].(string)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}
				_type, ok := claims["type"].(string)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}
				accountID, ok := claims["sub"].(string)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}
				issuer, ok := claims["iss"].(string)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}
				role, ok := claims["role"].(string)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}
				credential, ok := claims["credential"].(string)
				if !ok {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}

				accessToken := &AccessToken{
					Source:      source,
					Type:        _type,
					AccountID:   accountID,
					Credential:  credential,
					Issuer:      issuer,
					Role:        role,
					AccessToken: tokenString,
				}

				ctx := context.WithValue(r.Context(), ContextKeySource, source)
				ctx = context.WithValue(ctx, ContextKeyType, _type)
				ctx = context.WithValue(ctx, ContextKeyAccountID, accountID)
				ctx = context.WithValue(ctx, ContextKeyCredential, credential)
				ctx = context.WithValue(ctx, ContextKeyIssuer, issuer)
				ctx = context.WithValue(ctx, ContextKeyRole, role)
				ctx = context.WithValue(ctx, ContextKeyAccessToken, accessToken)

				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func GetAccessToken(ctx context.Context) (*AccessToken, bool) {
	accessToken, ok := ctx.Value(ContextKeyAccessToken).(*AccessToken)
	return accessToken, ok
}

func extractToken(r *http.Request) (string, bool) {
	if token, ok := extractTokenFromQueryParams(r); ok {
		return token, true
	}
	return extractTokenFromHeader(r)
}

func extractTokenFromQueryParams(r *http.Request) (string, bool) {
	token := r.URL.Query().Get("access_token")
	if token == "" {
		return "", false
	}
	return token, true
}

func extractTokenFromHeader(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(authHeader, "Bearer "), true
}

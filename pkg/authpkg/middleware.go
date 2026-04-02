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
	"go.uber.org/zap"
)

type JWTSecret interface {
	[]byte | *rsa.PublicKey | *rsa.PrivateKey
}

type Middleware func(http.Handler) http.Handler

type (
	contextKeyPermissions = struct{}
	contextKeyAccessToken = struct{}
)

func Protected(resource string, action Action, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		service := common.AppID
		permissions := GetPermissionsOrDefault(ctx)
		if !permissions.IsActionGranted(service, resource, action) {
			web.WriteBusinessErrorResponse(
				w, &common.ClientResponseError{
					HTTPCode: http.StatusForbidden,
					Message:  fmt.Sprintf("user does not have permission to %s %s", action, resource),
				},
			)
			return
		}

		next(w, r)
	}
}

func NewJWTMiddleware() func(http.Handler) http.Handler {
	secret := GetPublicKey()
	return JWTMiddleware[*rsa.PublicKey](jwt.SigningMethodRS256, secret, false)
}

func NewAdminJWTMiddleware() func(http.Handler) http.Handler {
	secret := GetAdminPublicKey()
	return JWTMiddleware[*rsa.PublicKey](jwt.SigningMethodRS256, secret, true)
}

func JWTMiddleware[S JWTSecret](method jwt.SigningMethod, secret S, admin bool) func(http.Handler) http.Handler {
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

				accessToken, err := ParseJWToken[S](
					method,
					secret,
					tokenString,
					admin,
				)
				if err != nil {
					web.WriteBusinessErrorResponse(
						w, &common.ClientResponseError{
							HTTPCode: http.StatusUnauthorized,
							Message:  "invalid token",
						},
					)
					return
				}

				ctx := context.WithValue(r.Context(), contextKeyPermissions{}, accessToken.Permissions)
				ctx = context.WithValue(ctx, contextKeyAccessToken{}, accessToken)

				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func ParseJWToken[S JWTSecret](
	method jwt.SigningMethod,
	secret S,
	tokenString string,
	admin bool,
) (*AccessToken, error) {
	token, err := jwt.Parse(
		tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != method.Alg() {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		},
	)

	logger := common.GetLogger()
	if err != nil {
		logger.Warn("[access_token] error parsing token", zap.Bool("admin", admin), zap.Error(err))
		return nil, err
	}
	if !token.Valid {
		logger.Warn("[access_token] is not valid", zap.Bool("admin", admin))
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Warn("[access_token] claims are not valid", zap.Bool("admin", admin))
		return nil, err
	}

	source, ok := claims["source"].(string)
	if !ok {
		logger.Warn(
			"[access_token] 'source' not found in claims",
			zap.String("role_required", AdminAccountRole),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		return nil, err
	}
	_type, ok := claims["type"].(string)
	if !ok {
		logger.Warn(
			"[access_token] 'type' not found in claims",
			zap.String("role_required", AdminAccountRole),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		return nil, err
	}
	accountID, ok := claims["sub"].(string)
	if !ok {
		logger.Warn(
			"[access_token] 'sub' not found in claims",
			zap.String("role_required", AdminAccountRole),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		return nil, err
	}
	issuer, ok := claims["iss"].(string)
	if !ok {
		logger.Warn(
			"[access_token] 'iss' not found in claims",
			zap.String("role_required", AdminAccountRole),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		return nil, err
	}
	role, ok := claims["role"].(string)
	if !ok {
		logger.Warn(
			"[access_token] 'role' not found in claims",
			zap.String("role_required", AdminAccountRole),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		return nil, err
	}

	if admin && role != AdminAccountRole {
		logger.Warn(
			"[access_token] invalid role requested",
			zap.String("role_required", AdminAccountRole),
			zap.String("role_requested", role),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		return nil, err
	}

	credential, ok := claims["credential"].(string)
	if !ok {
		logger.Warn(
			"[access_token] 'credential' not found in claims",
			zap.String("role_required", AdminAccountRole),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		return nil, err
	}
	scope, ok := claims["scope"].(string)
	if !ok {
		lvl := "debug"
		if admin {
			lvl = "warn"
		}
		logger.L(
			lvl,
			"[access_token] 'scope' not found in claims",
			zap.String("role_required", AdminAccountRole),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
		)
		scope = ""
	}
	permissions, err := ParsePermissions(scope)
	if err != nil {
		logger.Error(
			"[access_token] error parsing scopes",
			zap.String("scope", scope),
			zap.Bool("admin", admin),
			zap.Any("claims", claims),
			zap.Error(err),
		)
		return nil, err
	}

	accessToken := &AccessToken{
		Source:      source,
		Type:        _type,
		AccountID:   accountID,
		Credential:  credential,
		Issuer:      issuer,
		Role:        role,
		Permissions: permissions,
		AccessToken: tokenString,
	}
	return accessToken, nil
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

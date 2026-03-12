package common

import (
	"net/http"
	"slices"

	"go.uber.org/zap"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

var httpMethods = [9]string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func NewCors(logger Logger) *CORSConfig {
	allowedOrigins := GetEnvStringSlice("APP_CORS_ALLOWED_ORIGINS", []string{"*"})
	allowedMethods := GetEnvStringSlice("APP_CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := GetEnvStringSlice("APP_CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Authorization"})
	allowCredentials := GetEnvBool("APP_CORS_ALLOW_CREDENTIALS", false)

	for _, method := range allowedMethods {
		if !slices.Contains(httpMethods[:], method) {
			logger.Fatal("Method not allowed", zap.String("method", method))
		}
	}

	return &CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: allowCredentials,
	}
}

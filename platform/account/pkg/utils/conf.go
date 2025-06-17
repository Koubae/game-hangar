package utils

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"path/filepath"
)

func GetPrivateKey() (*rsa.PrivateKey, error) {
	filePath := filepath.Join("conf", "cert_private.pem")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(data)

}

func GetPublicKey() (*rsa.PublicKey, error) {
	filePath := filepath.Join("conf", "cert_public.pem")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(data)
}

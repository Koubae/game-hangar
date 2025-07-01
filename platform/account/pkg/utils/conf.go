package utils

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"path/filepath"
)

func GetPrivateKey() (*rsa.PrivateKey, error) {
	confDirName := GetEnvString("APP_CONF_DIR_NAME", "conf")

	filePath := filepath.Join(confDirName, "cert_private.pem")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(data)

}

func GetPublicKey() (*rsa.PublicKey, error) {
	confDirName := GetEnvString("APP_CONF_DIR_NAME", "conf")

	filePath := filepath.Join(confDirName, "cert_public.pem")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(data)
}

func GetPublicKeyOrPanic() *rsa.PublicKey {
	publicKey, err := GetPublicKey()
	if err != nil {
		panic(err.Error())
	}
	return publicKey
}

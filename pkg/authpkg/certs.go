package authpkg

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
	vars "github.com/koubae/game-hangar"
	"github.com/koubae/game-hangar/pkg/common"
)

var publicKey *rsa.PublicKey
var privateKey *rsa.PrivateKey

func LoadCerts(envPrefix string) error {

	publicFilePath := common.GetEnvString(
		envPrefix+"AUTH_PUBLIC_KEY_PATH",
		filepath.Join(vars.ConfigDir, "cert_public.pem"),
	)
	if !filepath.IsAbs(publicFilePath) {
		publicFilePath = filepath.Join(vars.RootDir, publicFilePath)
	}

	privateFilePath := common.GetEnvString(
		envPrefix+"AUTH_PRIVATE_KEY_PATH",
		filepath.Join(vars.ConfigDir, "cert_private.pem"),
	)
	if !filepath.IsAbs(privateFilePath) {
		privateFilePath = filepath.Join(vars.RootDir, privateFilePath)
	}

	public, err := os.ReadFile(publicFilePath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}
	private, err := os.ReadFile(privateFilePath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	return nil
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

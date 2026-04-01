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

var publicAdminKey *rsa.PublicKey
var privateAdminKey *rsa.PrivateKey

func LoadCerts(envPrefix string) error {
	// -------------------------------------------------------------------------
	// NOTE: Normal Certs
	// -------------------------------------------------------------------------
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

	// -------------------------------------------------------------------------
	// NOTE: Admin Certs
	// -------------------------------------------------------------------------
	adminPublicFilePath := common.GetEnvString(
		envPrefix+"AUTH_ADMIN_PUBLIC_KEY_PATH",
		filepath.Join(vars.ConfigDir, "cert_admin_public.pem"),
	)
	if !filepath.IsAbs(adminPublicFilePath) {
		adminPublicFilePath = filepath.Join(vars.RootDir, adminPublicFilePath)
	}

	adminPrivateFilePath := common.GetEnvString(
		envPrefix+"AUTH_ADMIN_PRIVATE_KEY_PATH",
		filepath.Join(vars.ConfigDir, "cert_admin_private.pem"),
	)
	if !filepath.IsAbs(adminPrivateFilePath) {
		adminPrivateFilePath = filepath.Join(vars.RootDir, adminPrivateFilePath)
	}

	adminPublic, err := os.ReadFile(adminPublicFilePath)
	if err != nil {
		return fmt.Errorf("failed to read Admin public key: %w", err)
	}
	adminPrivate, err := os.ReadFile(adminPrivateFilePath)
	if err != nil {
		return fmt.Errorf("failed to read Admin private key: %w", err)
	}

	publicAdminKey, err = jwt.ParseRSAPublicKeyFromPEM(adminPublic)
	if err != nil {
		return fmt.Errorf("failed to parse Admin public key: %w", err)
	}
	privateAdminKey, err = jwt.ParseRSAPrivateKeyFromPEM(adminPrivate)
	if err != nil {
		return fmt.Errorf("failed to parse Admin private key: %w", err)
	}

	return nil
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

func GetAdminPublicKey() *rsa.PublicKey {
	return publicAdminKey
}

func GetAdminPrivateKey() *rsa.PrivateKey {
	return privateAdminKey
}

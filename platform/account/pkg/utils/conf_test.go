package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/koubae/game-hangar/account/pkg/testings"
)

func TestGetPrivateKey(t *testing.T) {
	confDirName := GetEnvString("APP_CONF_DIR_NAME", "unit_conf")

	err := os.MkdirAll(confDirName, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	defer os.RemoveAll(confDirName)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyPath := filepath.Join(confDirName, "cert_private.pem")
	err = os.WriteFile(privateKeyPath, pem.EncodeToMemory(privateKeyPEM), 0600)
	if err != nil {
		t.Fatalf("Failed to write private key file: %v", err)
	}

	tests := []struct {
		name    string
		wantErr bool
		setup   func()
	}{
		{
			name:    "Valid private key file",
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Missing private key file",
			wantErr: true,
			setup: func() {
				os.Remove(privateKeyPath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := GetPrivateKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("GetPrivateKey() returned nil, want valid key")
			}
		})
	}
}

func TestGetPublicKey(t *testing.T) {
	confDirName := GetEnvString("APP_CONF_DIR_NAME", "unit_conf")
	err := os.MkdirAll(confDirName, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(confDirName)

	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Encode public key to PEM
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	// Write public key to file
	publicKeyPath := filepath.Join(confDirName, "cert_public.pem")
	err = os.WriteFile(publicKeyPath, pem.EncodeToMemory(publicKeyPEM), 0644)
	if err != nil {
		t.Fatalf("Failed to write public key file: %v", err)
	}

	tests := []struct {
		name    string
		wantErr bool
		setup   func()
	}{
		{
			name:    "Valid public key file",
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Missing public key file",
			wantErr: true,
			setup: func() {
				os.Remove(publicKeyPath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := GetPublicKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("GetPublicKey() returned nil, want valid key")
			}
		})
	}
}

func TestGetPublicKeyOrPanic(t *testing.T) {
	confDirName := GetEnvString("APP_CONF_DIR_NAME", "unit_conf")
	err := os.MkdirAll(confDirName, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(confDirName)

	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Encode public key to PEM
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	publicKeyPath := filepath.Join(confDirName, "cert_public.pem")

	tests := []struct {
		name        string
		shouldPanic bool
		setup       func()
	}{
		{
			name:        "Valid public key file",
			shouldPanic: false,
			setup: func() {
				err := os.WriteFile(publicKeyPath, pem.EncodeToMemory(publicKeyPEM), 0644)
				if err != nil {
					t.Fatalf("Failed to write public key file: %v", err)
				}
			},
		},
		{
			name:        "Missing public key file",
			shouldPanic: true,
			setup: func() {
				os.Remove(publicKeyPath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			defer func() {
				r := recover()
				if (r != nil) != tt.shouldPanic {
					t.Errorf("GetPublicKeyOrPanic() panic = %v, shouldPanic %v", r, tt.shouldPanic)
				}
			}()

			got := GetPublicKeyOrPanic()
			if !tt.shouldPanic && got == nil {
				t.Error("GetPublicKeyOrPanic() returned nil, want valid key")
			}
		})
	}
}

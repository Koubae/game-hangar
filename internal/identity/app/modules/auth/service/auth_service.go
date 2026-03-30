package service

import (
	"github.com/koubae/game-hangar/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

type AuthServiceFactory func() *AuthService

func NewAuthService() *AuthService {
	s := &AuthService{}
	return s
}

func (s *AuthService) HashSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.Wrap(errs.AuthSecretHash, err)
	}

	return string(hash), nil
}

func (s *AuthService) VerifySecret(hash string, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}

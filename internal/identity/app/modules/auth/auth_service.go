package auth

import (
	"github.com/koubae/game-hangar/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type SecretsService struct{}

type SecretsServiceFactory func() *SecretsService

func NewSecretsService() *SecretsService {
	s := &SecretsService{}
	return s
}

func (s *SecretsService) HashSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.Wrap(errs.AuthSecretHash, err)
	}

	return string(hash), nil
}

func (s *SecretsService) VerifySecret(hash string, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}

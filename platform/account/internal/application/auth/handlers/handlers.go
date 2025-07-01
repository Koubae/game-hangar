package handlers

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/koubae/game-hangar/account/internal/settings"
	"github.com/koubae/game-hangar/account/pkg/utils"
	"strings"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
}

func (r *LoginRequest) Validate() error {
	// Normalize Data
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)
	r.ClientID = strings.TrimSpace(r.ClientID)

	if r.Username == "" {
		return errors.New("username is required")
	} else if r.Password == "" {
		return errors.New("password is required")
	} else if r.ClientID == "" {
		return errors.New("client_id is required")
	}

	return nil
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

type LoginHandler struct {
	Command  LoginRequest
	Response LoginResponse
}

func (h *LoginHandler) Handle() error {
	expire := time.Now().Add(settings.AuthTokenExpirationTime).Unix()
	token, err := generateJWTWithRSA(h.Command.Username, h.Command.ClientID, expire)
	if err != nil {
		return err
	}

	h.Response = LoginResponse{AccessToken: token, Expires: expire}
	return nil
}

func generateJWTWithRSA(userID string, clientID string, expire int64) (string, error) {
	privateKey := loadAndGetPrivateKey()
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expire,
		"iss": "game-hangar", // could be dynamic? services would send their identifier and other have default

		"role":      "user",
		"client_id": clientID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)

}

var privateKey *rsa.PrivateKey

func loadAndGetPrivateKey() *rsa.PrivateKey {
	var err error
	if privateKey == nil {
		privateKey, err = utils.GetPrivateKey()
		if err != nil {
			panic(err.Error())
		}
	}
	return privateKey
}

package controllers

import (
	"crypto/rsa"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/koubae/game-hangar/account/pkg/utils"
	"net/http"
	"strings"
	"time"
)

type AuthController struct{}

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

func (controller *AuthController) LoginV1(c *gin.Context) {
	var request = LoginRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expire := time.Now().Add(AuthTokenExpirationTime).Unix()
	token, err := GenerateJWTWithHMAC(request.Username, request.ClientID, expire)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := LoginResponse{AccessToken: token, Expires: expire}
	c.JSON(http.StatusOK, response)
}

const AuthTokenExpirationTime = time.Hour * 4

// TODO: .env variable!
var AUTH_SECRET = []byte("AUTH_SECRET_1234")

//var AUTH_SECRET = []byte(os.Getenv("AUTH_JWT_SECRET"))

// TODO: Move into a service!
func GenerateJWTWithHMAC(userID string, clientID string, expire int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expire,
		"iss": "game-hangar", // could be dynamic? services would send their identifier and other have default

		"role":      "user",
		"client_id": clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AUTH_SECRET)
}

func (controller *AuthController) LoginV2(c *gin.Context) {
	var request = LoginRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expire := time.Now().Add(AuthTokenExpirationTime).Unix()
	token, err := GenerateJWTWithRSA(request.Username, request.ClientID, expire)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := LoginResponse{AccessToken: token, Expires: expire}
	c.JSON(http.StatusOK, response)
}

func GenerateJWTWithRSA(userID string, clientID string, expire int64) (string, error) {
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

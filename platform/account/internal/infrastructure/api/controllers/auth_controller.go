package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/koubae/game-hangar/account/internal/application/auth/handlers"
	"net/http"
)

type AuthController struct{}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

func (controller *AuthController) LoginV1(c *gin.Context) {
	var request = handlers.LoginRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler := handlers.LoginHandler{Command: request}

	err := handler.Handle()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)
}

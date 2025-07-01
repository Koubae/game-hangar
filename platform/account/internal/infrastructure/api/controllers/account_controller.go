package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/koubae/game-hangar/account/internal/application/account/handlers"
)

type AccountControllers struct{}

func (controller *AccountControllers) Get(c *gin.Context) {
	var request = handlers.GetAccountRequest{
		Username: c.Params.ByName("name"),
		ClientID: c.MustGet("client_id").(string),
		UserID:   c.MustGet("user_id").(uint),
	}

	err := c.Bind(&request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	handler := handlers.GetAccountHandler{Command: request}
	if err = handler.Handle(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}) // TODO: check error type!
		return
	}

	c.JSON(200, handler.Response)
}

package controllers

import "github.com/gin-gonic/gin"

type AccountControllers struct{}

type GetRequest struct {
	FullProfile bool `form:"full_profile" json:"full_profile"`
}

func (controller *AccountControllers) Get(c *gin.Context) {
	username := c.Params.ByName("name")

	request := GetRequest{}
	err := c.Bind(&request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// TODO should be able to get admin/no admin roles?
	clientID := c.MustGet("client_id").(string)
	userID := c.MustGet("user_id").(string)

	// TODO :)
	c.JSON(200, gin.H{"username": username, "request": request, "client_id": clientID, "user_id": userID})
}

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Account struct {
	Username string `json:"username"`
	ClientID string `json:"client_id"`
}

func (controller *AccountControllers) Create(c *gin.Context) {
	// TODO: Add body read
	var request = CreateRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	clientID := c.MustGet("client_id").(string)

	account := Account{Username: request.Username, ClientID: clientID}
	c.JSON(200, account)
}

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

	// TODO :)
	c.JSON(200, gin.H{"username": username, "request": request})
}

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Account struct {
	Username string `json:"username"`
}

func (controller *AccountControllers) Create(c *gin.Context) {
	// TODO: Add body read
	var request = CreateRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	account := Account{Username: request.Username}
	c.JSON(200, account)
}

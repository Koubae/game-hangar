package controllers

import "github.com/gin-gonic/gin"

func (controller *AccountControllers) Get(c *gin.Context) {
	username := c.Params.ByName("name")

	// TODO :)
	c.JSON(200, gin.H{"username": username})
}

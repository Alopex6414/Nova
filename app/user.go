package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleCreatUser(c *gin.Context) {
	var user User
	// request body should bind json
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, user)
	return
}

func handleDeleteUser(c *gin.Context) {
	return
}

func handleModifyUser(c *gin.Context) {
	return
}

func handleQueryUser(c *gin.Context) {
	return
}

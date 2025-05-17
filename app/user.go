package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func HandleCreateUserId(c *gin.Context) {
	// create userId
	var userId UserID
	userId.UserId = uuid.New().String()
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, userId)
}

func HandleCreateUser(c *gin.Context) {
	var request User
	// request body should bind json
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// return response
	var response = request
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
	return
}

func HandleDeleteUser(c *gin.Context) {
	return
}

func HandleModifyUser(c *gin.Context) {
	return
}

func HandleQueryUser(c *gin.Context) {
	return
}

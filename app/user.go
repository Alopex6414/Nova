package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (nova *Nova) HandleCreateUserId(c *gin.Context) {
	// create userId
	var userId UserID
	userId.UserId = uuid.New().String()
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, userId)
}

func (nova *Nova) HandleCreateUser(c *gin.Context) {
	var request User
	// request body should bind json
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check user existence
	b := func(userId string) bool {
		// enable user cache read lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search userId in data cache
		for _, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				return true
			}
		}
		return false
	}(request.UserId)
	if b {
		var problemDetails ProblemDetails
		problemDetails.Title = "User Conflict"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusConflict
		problemDetails.Cause = errors.New("user already exists").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusConflict, problemDetails)
		return
	}
	// store created user in data cache
	response := User{
		UserId:      strings.ToLower(request.UserId),
		Username:    request.Username,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Address:     request.Address,
		Company:     request.Company,
	}
	func(user User) {
		// enable user cache write lock
		nova.cache.userCache.mutex.Lock()
		defer nova.cache.userCache.mutex.Unlock()
		// append user in data cache
		nova.cache.userCache.userSet = append(nova.cache.userCache.userSet, user)
	}(response)
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
	return
}

func (nova *Nova) HandleDeleteUser(c *gin.Context) {
	// extract userId from uri
	userId := strings.ToLower(c.Param("userId"))
	// check user existence
	b := func(userId string) bool {
		// enable user cache read lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search userId in data cache
		for _, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				return true
			}
		}
		return false
	}(userId)
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "User Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("user not found").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// delete user from data cache
	func(userId string) {
		// enable user cache write lock
		nova.cache.userCache.mutex.Lock()
		defer nova.cache.userCache.mutex.Unlock()
		// search & delete user from data cache
		for k, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				nova.cache.userCache.userSet = append(nova.cache.userCache.userSet[:k], nova.cache.userCache.userSet[k+1:]...)
				break
			}
		}
		return
	}(userId)
	// return response
	c.Status(http.StatusNoContent)
	return
}

func (nova *Nova) HandleModifyUser(c *gin.Context) {
	return
}

func (nova *Nova) HandleQueryUser(c *gin.Context) {
	return
}

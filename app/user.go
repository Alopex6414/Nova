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
	c.JSON(http.StatusCreated, userId)
	return
}

func (nova *Nova) HandleQueryUserId(c *gin.Context) {
	var userName UserName
	// request body should bind json
	err := c.ShouldBindJSON(&userName)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// query userId from data cache
	userId, err := func(userName UserName) (UserID, error) {
		// enable user cache write lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search & delete user from data cache
		for k, v := range nova.cache.userCache.userSet {
			if v.Username == userName.Username {
				return UserID{UserId: nova.cache.userCache.userSet[k].UserId}, nil
			}
		}
		return UserID{}, errors.New("userId not found")
	}(userName)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("userId not found").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, userId)
	return
}

func (nova *Nova) HandleCreateUser(c *gin.Context) {
	var request User
	// request body should bind json
	err := c.ShouldBindJSON(&request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// check request body correctness
	b, err := func(user User) (bool, error) {
		// check userId format is UUID
		err = uuid.Validate(user.UserId)
		if err != nil {
			return false, err
		}
		return true, nil
	}(request)
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// check user existence
	b = func(userId string) bool {
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
	}(strings.ToLower(request.UserId))
	if b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Conflict"
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
	// store created user in database
	err = func(userId string) error {
		// enable user cache read lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search userId in data cache
		b := false
		user := User{}
		for _, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				user = v
				b = true
				break
			}
		}
		if !b {
			return errors.New("user not found")
		}
		// create user in database
		_, err := nova.db.CreateUser(&user)
		if err != nil {
			return err
		}
		return nil
	}(response.UserId)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
	return
}

func (nova *Nova) HandleDeleteUser(c *gin.Context) {
	// extract userId from uri
	userId := strings.ToLower(c.Param("userId"))
	// request userId correctness
	b, _ := func(userId string) (bool, error) {
		// check userId format is UUID
		err := uuid.Validate(userId)
		if err != nil {
			return false, err
		}
		return true, nil
	}(userId)
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = errors.New("userId format incorrect").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// check user existence
	b = func(userId string) bool {
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
		problemDetails.Title = "Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("user not found").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// delete user from database
	err := func(userId string) error {
		// enable user cache read lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search userId in data cache
		b := false
		for _, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				b = true
				break
			}
		}
		if !b {
			return errors.New("user not found")
		}
		// delete user in database
		err := nova.db.DeleteUser(userId)
		if err != nil {
			return err
		}
		return nil
	}(userId)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusInternalServerError, problemDetails)
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
	var request User
	// request body should bind json
	err := c.ShouldBindJSON(&request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// check request body correctness
	b, err := func(user User) (bool, error) {
		// check userId format is UUID
		err = uuid.Validate(user.UserId)
		if err != nil {
			return false, err
		}
		return true, nil
	}(request)
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// check user existence
	b = func(userId string) bool {
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
	}(strings.ToLower(request.UserId))
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("user not found").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// store modified user in data cache
	response, err := func(user User) (User, error) {
		// enable user cache write lock
		nova.cache.userCache.mutex.Lock()
		defer nova.cache.userCache.mutex.Unlock()
		// replace user in data cache
		for k, v := range nova.cache.userCache.userSet {
			if v.UserId == user.UserId {
				if user.Username != "" {
					nova.cache.userCache.userSet[k].Username = user.Username
				}
				if user.Password != "" {
					nova.cache.userCache.userSet[k].Password = user.Password
				}
				if user.PhoneNumber != "" {
					nova.cache.userCache.userSet[k].PhoneNumber = user.PhoneNumber
				}
				if user.Email != "" {
					nova.cache.userCache.userSet[k].Email = user.Email
				}
				if user.Address != "" {
					nova.cache.userCache.userSet[k].Address = user.Address
				}
				if user.Company != "" {
					nova.cache.userCache.userSet[k].Company = user.Company
				}
				return nova.cache.userCache.userSet[k], nil
			}
		}
		return User{}, errors.New("user not found")
	}(request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("user not found").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// store patched user in database
	err = func(userId string) error {
		// enable user cache read lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search user in data cache
		b := false
		user := User{}
		for _, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				user = v
				b = true
				break
			}
		}
		if !b {
			return errors.New("user not found")
		}
		// update user in database
		err = nova.db.UpdateUser(&user)
		if err != nil {
			return err
		}
		return nil
	}(response.UserId)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
	return
}

func (nova *Nova) HandleQueryUser(c *gin.Context) {
	// extract userId from uri
	userId := strings.ToLower(c.Param("userId"))
	// request userId correctness
	b, _ := func(userId string) (bool, error) {
		// check userId format is UUID
		err := uuid.Validate(userId)
		if err != nil {
			return false, err
		}
		return true, nil
	}(userId)
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = errors.New("userId format incorrect").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// check user existence
	b = func(userId string) bool {
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
		problemDetails.Title = "Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("user not found").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// query user from database
	err := func(userId string) error {
		// enable user cache write lock
		nova.cache.userCache.mutex.Lock()
		defer nova.cache.userCache.mutex.Unlock()
		// query user from database
		user, err := nova.db.QueryUser(userId)
		if err != nil {
			return err
		}
		// update user in data cache
		for k, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				nova.cache.userCache.userSet[k] = *user
			}
		}
		return nil
	}(userId)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}
	// query user from data cache
	response, err := func(userId string) (User, error) {
		// enable user cache read lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search & query user from data cache
		for k, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				return nova.cache.userCache.userSet[k], nil
			}
		}
		return User{}, errors.New("user not found")
	}(userId)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("user not found").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
	return
}

func (nova *Nova) HandleUpdateUser(c *gin.Context) {
	var request User
	// request body should bind json
	err := c.ShouldBindJSON(&request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	// check request body correctness
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
	}(strings.ToLower(request.UserId))
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Forbidden"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusForbidden
		problemDetails.Cause = errors.New("forbidden replace user without create it").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusForbidden, problemDetails)
		return
	}
	// store updated user in data cache
	response := User{
		UserId:      strings.ToLower(request.UserId),
		Username:    request.Username,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Address:     request.Address,
		Company:     request.Company,
	}
	b = func(user User) bool {
		// enable user cache write lock
		nova.cache.userCache.mutex.Lock()
		defer nova.cache.userCache.mutex.Unlock()
		// replace user in data cache
		for k, v := range nova.cache.userCache.userSet {
			if v.UserId == user.UserId {
				nova.cache.userCache.userSet[k] = user
				return true
			}
		}
		return false
	}(response)
	if !b {
		var problemDetails ProblemDetails
		problemDetails.Title = "Not Found"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Cause = errors.New("user not existed").Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}
	// store update user in database
	err = func(userId string) error {
		// enable user cache read lock
		nova.cache.userCache.mutex.RLock()
		defer nova.cache.userCache.mutex.RUnlock()
		// search user in data cache
		b := false
		user := User{}
		for _, v := range nova.cache.userCache.userSet {
			if v.UserId == userId {
				user = v
				b = true
				break
			}
		}
		if !b {
			return errors.New("user not found")
		}
		// update user in database
		err = nova.db.UpdateUser(&user)
		if err != nil {
			return err
		}
		return nil
	}(response.UserId)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Type = "User"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Cause = err.Error()
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}
	// return response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
	return
}

func (nova *Nova) HandleCreateUserLogin(c *gin.Context) {
	return
}

package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nova/logger"
	"strings"
)

func (nova *Nova) HandleCreateUserId(c *gin.Context) {
	// create userId
	var userId UserID
	logger.Infof("handle request create userId")
	// generate userId
	userId.UserId = uuid.New().String()
	logger.Debugf("generate userId: %v", userId.UserId)
	// return response
	nova.response201Created(c, userId)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, userId.UserId)
	return
}

func (nova *Nova) HandleQueryUserId(c *gin.Context) {
	// query userId
	var userName UserName
	logger.Infof("handle request query userId")
	// request body should bind json
	logger.Debugf("request body bind json format")
	err := c.ShouldBindJSON(&userName)
	if err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// query userId from data cache
	logger.Debugf("query userId from data cache")
	userId, err := nova.queryUserFromDataCache(userName)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error query user from data cache: %v", err)
		return
	}
	logger.Debugf("successfully query user from data cache")
	// return response
	nova.response200OK(c, userId)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, userId.UserId)
	return
}

func (nova *Nova) HandleCreateUser(c *gin.Context) {
	// create user
	var request User
	logger.Infof("handle request create user")
	// request body should bind json
	logger.Debugf("request body bind json format")
	err := c.ShouldBindJSON(&request)
	if err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check user is validate")
	b, err := nova.isUserValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check user is validate: %v", err)
		return
	}
	logger.Debugf("successfully check user is validate")
	// check user existence
	logger.Debugf("check user is existed")
	if nova.isUserExisted(strings.ToLower(request.UserId)) {
		nova.response409Conflict(c, errors.New("user already exists"))
		logger.Errorf("error check user is existed: %v", err)
		return
	}
	logger.Debugf("successfully check user is existed")
	// store created user in data cache
	logger.Debugf("store user in data cache")
	response := User{
		UserId:      strings.ToLower(request.UserId),
		Username:    request.Username,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Address:     request.Address,
		Company:     request.Company,
	}
	nova.createUserInDataCache(response)
	logger.Debugf("successfully store user in data cache")
	// store created user in database
	logger.Debugf("store user in database")
	if err = nova.createUserInDatabase(response.UserId); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error user in database: %v", err)
		return
	}
	logger.Debugf("successfully store user in database")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) HandleDeleteUser(c *gin.Context) {
	// delete user
	logger.Infof("handle request delete user")
	// extract userId from uri
	userId := strings.ToLower(c.Param("userId"))
	// request userId correctness
	logger.Debugf("check userId is validate")
	if b, _ := nova.isUserIdValidate(userId); !b {
		nova.response400BadRequest(c, errors.New("userId format incorrect"))
		logger.Error("error check userId is validate")
		return
	}
	logger.Debugf("successfully check userId is validate")
	// check user existence
	logger.Debugf("check user is validate")
	if !nova.isUserExisted(userId) {
		nova.response404NotFound(c, errors.New("user not found"))
		logger.Error("error check user is validate")
		return
	}
	logger.Debugf("successfully check user is validate")
	// delete user from database
	logger.Debugf("delete user in database")
	if err := nova.deleteUserInDatabase(userId); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Error("error delete user in database")
		return
	}
	logger.Debugf("successfully delete user in database")
	// delete user from data cache
	logger.Debugf("delete user in data cache")
	nova.deleteUserInDataCache(userId)
	// return response
	nova.response204NoContent(c, nil)
	logger.Infof("response status code: %v, body: %v", http.StatusNoContent, nil)
	return
}

func (nova *Nova) HandleModifyUser(c *gin.Context) {
	// modify user
	var request User
	logger.Infof("handle request modify user")
	// request body should bind json
	logger.Debugf("request body bind json format")
	err := c.ShouldBindJSON(&request)
	if err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check user is validate")
	b, err := nova.isUserValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check user is existed: %v", err)
		return
	}
	logger.Debugf("successfully check user is validate")
	// check user existence
	logger.Debugf("check user is existed")
	if !nova.isUserExisted(strings.ToLower(request.UserId)) {
		nova.response404NotFound(c, errors.New("user not found"))
		logger.Errorf("error check user is existed: %v", err)
		return
	}
	logger.Debugf("successfully check user is existed")
	// store modified user in data cache
	logger.Debugf("store modify user in data cache")
	response, err := nova.modifyUserInDataCache(request)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error store modify user in data cache: %v", err)
		return
	}
	logger.Debugf("succefully store modify user in data cache")
	// store patched user in database
	logger.Debugf("store modify user in database")
	if err = nova.modifyUserInDatabase(response.UserId); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error store modify user in database: %v", err)
		return
	}
	logger.Debugf("successfully store modify user in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleQueryUser(c *gin.Context) {
	// query user
	logger.Infof("handle request query user")
	// extract userId from uri
	userId := strings.ToLower(c.Param("userId"))
	// request userId correctness
	logger.Debugf("check userId is validate")
	if b, _ := nova.isUserIdValidate(userId); !b {
		nova.response400BadRequest(c, errors.New("userId format incorrect"))
		logger.Errorf("error check userId is validate")
		return
	}
	logger.Debugf("successfully check userId is validate")
	// check user existence
	logger.Debugf("check user is existed")
	if !nova.isUserExisted(userId) {
		nova.response404NotFound(c, errors.New("user not found"))
		logger.Errorf("error check user is existed")
		return
	}
	logger.Debugf("successfully check user is existed")
	// query user from database
	logger.Debugf("query user in database")
	if err := nova.queryUserInDatabase(userId); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error query user in database: %v", err)
		return
	}
	logger.Debugf("successfully query user in database")
	// query user from data cache
	logger.Debugf("query user in data cache")
	response, err := nova.queryUserInDataCache(userId)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error query user in data cache: %v", err)
		return
	}
	logger.Debugf("successfully query user in data cache")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleUpdateUser(c *gin.Context) {
	var request User
	// request body should bind json
	err := c.ShouldBindJSON(&request)
	if err != nil {
		nova.response400BadRequest(c, err)
		return
	}
	// check request body correctness
	// check user existence
	if !nova.isUserExisted(strings.ToLower(request.UserId)) {
		nova.response403Forbidden(c, errors.New("forbidden replace user without create it"))
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
	b := func(user User) bool {
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
		nova.response404NotFound(c, errors.New("user not found"))
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
		nova.response500InternalServerError(c, err)
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

func (nova *Nova) isUserExisted(userId string) bool {
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
}

func (nova *Nova) isUserIdValidate(userId string) (bool, error) {
	// check userId format is UUID
	err := uuid.Validate(userId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) isUserValidate(user User) (bool, error) {
	// check userId format is UUID
	err := uuid.Validate(user.UserId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) createUserInDataCache(user User) {
	// enable user cache write lock
	nova.cache.userCache.mutex.Lock()
	defer nova.cache.userCache.mutex.Unlock()
	// append user in data cache
	nova.cache.userCache.userSet = append(nova.cache.userCache.userSet, user)
	return
}

func (nova *Nova) deleteUserInDataCache(userId string) {
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
}

func (nova *Nova) modifyUserInDataCache(user User) (User, error) {
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
}

func (nova *Nova) queryUserInDataCache(userId string) (User, error) {
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
}

func (nova *Nova) createUserInDatabase(userId string) error {
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
	if _, err := nova.db.CreateUser(&user); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) deleteUserInDatabase(userId string) error {
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
	if err := nova.db.DeleteUser(userId); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) modifyUserInDatabase(userId string) error {
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
	if err := nova.db.UpdateUser(&user); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) queryUserInDatabase(userId string) error {
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
}

func (nova *Nova) queryUserFromDataCache(userName UserName) (UserID, error) {
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
}

func (nova *Nova) response200OK(c *gin.Context, body any) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, body)
	return
}

func (nova *Nova) response201Created(c *gin.Context, body any) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, body)
	return
}

func (nova *Nova) response204NoContent(c *gin.Context, body any) {
	c.Status(http.StatusNoContent)
	return
}

func (nova *Nova) response400BadRequest(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Bad Request"
	problemDetails.Type = "User"
	problemDetails.Status = http.StatusBadRequest
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusBadRequest, problemDetails)
	return
}

func (nova *Nova) response403Forbidden(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Forbidden"
	problemDetails.Type = "User"
	problemDetails.Status = http.StatusForbidden
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusForbidden, problemDetails)
	return
}

func (nova *Nova) response404NotFound(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Not Found"
	problemDetails.Type = "User"
	problemDetails.Status = http.StatusNotFound
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusNotFound, problemDetails)
	return
}

func (nova *Nova) response409Conflict(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Conflict"
	problemDetails.Type = "User"
	problemDetails.Status = http.StatusConflict
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusConflict, problemDetails)
	return
}

func (nova *Nova) response500InternalServerError(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Internal Server Error"
	problemDetails.Type = "User"
	problemDetails.Status = http.StatusInternalServerError
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusInternalServerError, problemDetails)
	return
}

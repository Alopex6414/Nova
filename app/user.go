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
	var userId string
	logger.Infof("handle request create userId")
	// generate userId
	userId = uuid.New().String()
	logger.Debugf("generate userId: %v", userId)
	// return response
	nova.response201Created(c, userId)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, userId)
	return
}

func (nova *Nova) HandleQueryUserId(c *gin.Context) {
	// query userId
	var userName string
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
	logger.Infof("response status code: %v, body: %v", http.StatusOK, userId)
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
	// update data cache by querying users in database
	logger.Debugf("update data cache by querying users in database")
	err = nova.queryUsersInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying users in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying users in database")
	// check user existence
	logger.Debugf("check user is existed")
	if nova.isUserExisted(strings.ToLower(request.UserId)) {
		nova.response409Conflict(c, errors.New("user already exists"))
		logger.Errorf("error check user is existed: %v", err)
		return
	}
	logger.Debugf("successfully check user is existed")
	// check userName or phoneNumber existence
	logger.Debugf("check userName or phoneNumber is existed")
	if nova.isUserNameOrPhoneExisted(request.Username, request.PhoneNumber) {
		nova.response409Conflict(c, errors.New("userName or phoneNumber already exists"))
		logger.Errorf("error check userName or phoneNumber is existed: %v", err)
		return
	}
	logger.Debugf("successfully check userName or phoneNumber is existed")
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
	// update user
	var request User
	logger.Infof("handle request update user")
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
	// check user existence
	logger.Debugf("check user existence")
	if !nova.isUserExisted(strings.ToLower(request.UserId)) {
		nova.response403Forbidden(c, errors.New("forbidden replace user without create it"))
		logger.Errorf("error check user existence")
		return
	}
	logger.Debugf("successfully check user existence")
	// store updated user in data cache
	logger.Debugf("update user in data cache")
	response := User{
		UserId:      strings.ToLower(request.UserId),
		Username:    request.Username,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Address:     request.Address,
		Company:     request.Company,
	}
	if b := nova.updateUserInDataCache(response); !b {
		nova.response404NotFound(c, errors.New("user not found"))
		logger.Errorf("error update user in data cache")
		return
	}
	logger.Debugf("successfully update user in data cache")
	// store update user in database
	logger.Debugf("update user in database")
	if err = nova.updateUserInDatabase(response.UserId); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update user in database")
		return
	}
	logger.Debugf("successfully update user in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleCreateUserLogin(c *gin.Context) {
	// login user
	var request UserLogin
	logger.Infof("handle request login user")
	// extract userId from uri
	userId := strings.ToLower(c.Param("userId"))
	// request body should bind json
	logger.Debugf("request body bind json format")
	err := c.ShouldBindJSON(&request)
	if err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
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
	user, err := nova.queryUserInDataCache(userId)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error query user in data cache: %v", err)
		return
	}
	logger.Debugf("successfully query user in data cache")
	// check username
	logger.Debugf("check username consistentance")
	if user.Username != request.Username {
		nova.response412PreconditionFailed(c, errors.New("request username inconsistent with database"))
		logger.Errorf("error check username consistentance.")
		return
	}
	logger.Debugf("successfully check username consistentance")
	// verify password correctness
	logger.Debugf("check password correctness")
	if user.Password != request.Password {
		nova.response417ExpectationFailed(c, errors.New("request password inconsistent with database"))
		logger.Errorf("error check password correctness.")
		return
	}
	logger.Debugf("successfully check password correctness")
	// return response
	nova.response200OK(c, nil)
	logger.Infof("response status code: %v", http.StatusOK)
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

func (nova *Nova) isUserNameOrPhoneExisted(userName string, phoneNumber string) bool {
	// enable user cache read lock
	nova.cache.userCache.mutex.RLock()
	defer nova.cache.userCache.mutex.RUnlock()
	// search userId in data cache
	for _, v := range nova.cache.userCache.userSet {
		if v.Username == userName || v.PhoneNumber == phoneNumber {
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

func (nova *Nova) updateUserInDataCache(user User) bool {
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
}

func (nova *Nova) queryUsersInDataCache() ([]User, error) {
	// enable user cache read lock
	nova.cache.userCache.mutex.RLock()
	defer nova.cache.userCache.mutex.RUnlock()
	// search & query users from data cache
	return nova.cache.userCache.userSet, nil
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
			break
		}
	}
	return nil
}

func (nova *Nova) updateUserInDatabase(userId string) error {
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

func (nova *Nova) queryUsersInDatabase() error {
	// enable user cache write lock
	nova.cache.userCache.mutex.Lock()
	defer nova.cache.userCache.mutex.Unlock()
	// query user from database
	users, err := nova.db.QueryUsers()
	if err != nil {
		return err
	}
	// update user in data cache
	for _, user := range users {
		b := false
		// update if user existed
		for k, v := range nova.cache.userCache.userSet {
			if v.UserId == user.UserId {
				nova.cache.userCache.userSet[k] = *user
				b = true
				break
			}
		}
		// create user if user not existed
		if !b {
			nova.cache.userCache.userSet = append(nova.cache.userCache.userSet)
		}
	}
	return nil
}

func (nova *Nova) queryUserFromDataCache(userName string) (string, error) {
	// enable user cache write lock
	nova.cache.userCache.mutex.RLock()
	defer nova.cache.userCache.mutex.RUnlock()
	// search & delete user from data cache
	for k, v := range nova.cache.userCache.userSet {
		if v.Username == userName {
			return nova.cache.userCache.userSet[k].UserId, nil
		}
	}
	return "", errors.New("userId not found")
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

func (nova *Nova) response412PreconditionFailed(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "PreconditionFailed"
	problemDetails.Type = "User"
	problemDetails.Status = http.StatusPreconditionFailed
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusPreconditionFailed, problemDetails)
	return
}

func (nova *Nova) response417ExpectationFailed(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "ExpectationFailed"
	problemDetails.Type = "User"
	problemDetails.Status = http.StatusExpectationFailed
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusExpectationFailed, problemDetails)
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

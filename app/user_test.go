package app

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRouter() *gin.Engine {
	// apply default Gin service
	router := gin.Default()
	// apply Gin logger & recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// create router group for nova
	novaService := router.Group("nova/v1")
	{
		novaService.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "hello Gin\n") })
		/* user management */
		// userId related
		novaService.POST("/user/userId", HandleCreateUserId)
		// user related
		novaService.PUT("/user/:userId", HandleCreateUser)
		novaService.DELETE("/user/:userId", HandleDeleteUser)
		novaService.PATCH("/user/:userId", HandleModifyUser)
		novaService.GET("/user/:userId", HandleQueryUser)
	}
	return router
}

func startTestService() (*httptest.Server, *gin.Engine) {
	router := setupTestRouter()
	return httptest.NewServer(router), router
}

func TestHandleCreateUserId(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestHandleCreateUserId
	// Test Purpose: Test HandleCreateUserId create userId
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	---------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	// request content
	url := server.URL + "/nova/v1/user/userId"
	// request create userId
	w := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(w, request)
	// return response
	var response UserID
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.NoError(t, uuid.Validate(response.UserId))
}

func TestHandleCreateUser(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestHandleCreateUser
	// Test Purpose: Test HandleCreateUser create user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using PUT method
	// 4. receive CreateUser response with user information by using 201 Created Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	/* create userId */
	// request content
	url := server.URL + "/nova/v1/user/userId"
	// request create userId
	wUserId := httptest.NewRecorder()
	reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wUserId, reqUserId)
	// return response
	var resUserId UserID
	err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusOK, wUserId.Code)
	assert.Equal(t, "application/json", wUserId.Header().Get("Content-Type"))
	assert.NoError(t, uuid.Validate(resUserId.UserId))
	/* create user */
	// request content
	url = server.URL + "/nova/v1/user"
	user := User{
		UserId:      resUserId.UserId,
		Username:    "alice",
		Password:    "123456",
		PhoneNumber: "+1412387",
		Email:       "alice@gmail.com",
		Address:     "No.5, Wall Street, New York, USA",
		Company:     "Apple Inc.",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Errorf("error marshal user: %v", err)
	}
	// request create user
	wUser := httptest.NewRecorder()
	reqUser, err := http.NewRequest(http.MethodPut, url+"/"+resUserId.UserId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqUser.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wUser, reqUser)
	// return response
	var resUser User
	err = json.Unmarshal(wUser.Body.Bytes(), &resUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wUser.Code)
	assert.Equal(t, "application/json", wUser.Header().Get("Content-Type"))
	assert.Equal(t, user.UserId, resUser.UserId)
	assert.Equal(t, user.Username, resUser.Username)
	assert.Equal(t, user.Password, resUser.Password)
	assert.Equal(t, user.PhoneNumber, resUser.PhoneNumber)
	assert.Equal(t, user.Email, resUser.Email)
	assert.Equal(t, user.Address, resUser.Address)
	assert.Equal(t, user.Company, resUser.Company)
}

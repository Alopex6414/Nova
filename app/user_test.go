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
	// create Nova instance
	nova := New()
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
		novaService.POST("/user/userId", nova.HandleCreateUserId)
		// user related
		novaService.POST("/user/:userId", nova.HandleCreateUser)
		novaService.PUT("/user/:userId", nova.HandleUpdateUser)
		novaService.DELETE("/user/:userId", nova.HandleDeleteUser)
		novaService.PATCH("/user/:userId", nova.HandleModifyUser)
		novaService.GET("/user/:userId", nova.HandleQueryUser)
	}
	return router
}

func startTestService() (*httptest.Server, *gin.Engine) {
	router := setupTestRouter()
	return httptest.NewServer(router), router
}

func TestNova_HandleCreateUserId(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateUserId
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

func BenchmarkNova_HandleCreateUserId(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateUserId
	// Test Purpose: Benchmark HandleCreateUserId create userId
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	---------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	// start benchmark test
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// request content
		url := server.URL + "/nova/v1/user/userId"
		// request create userId
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(w, request)
		// return response
		var response UserID
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusOK, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.NoError(b, uuid.Validate(response.UserId))
	}
}

func BenchmarkNova_HandleCreateUserIdParallel(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateUserIdParallel
	// Test Purpose: Benchmark HandleCreateUserId create userId
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	---------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// request content
			url := server.URL + "/nova/v1/user/userId"
			// request create userId
			w := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(w, request)
			// return response
			var response UserID
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusOK, w.Code)
			assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
			assert.NoError(b, uuid.Validate(response.UserId))
		}
	})
}

func TestNova_HandleCreateUser(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateUser
	// Test Purpose: Test HandleCreateUser create user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
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
	reqUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
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

func BenchmarkNova_HandleCreateUser(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateUser
	// Test Purpose: Benchmark HandleCreateUser create user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		/* create userId */
		// request content
		url := server.URL + "/nova/v1/user/userId"
		// request create userId
		wUserId := httptest.NewRecorder()
		reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wUserId, reqUserId)
		// return response
		var resUserId UserID
		err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusOK, wUserId.Code)
		assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
		assert.NoError(b, uuid.Validate(resUserId.UserId))
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
			b.Errorf("error marshal user: %v", err)
		}
		// request create user
		wUser := httptest.NewRecorder()
		reqUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqUser.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wUser, reqUser)
		// return response
		var resUser User
		err = json.Unmarshal(wUser.Body.Bytes(), &resUser)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wUser.Code)
		assert.Equal(b, "application/json", wUser.Header().Get("Content-Type"))
		assert.Equal(b, user.UserId, resUser.UserId)
		assert.Equal(b, user.Username, resUser.Username)
		assert.Equal(b, user.Password, resUser.Password)
		assert.Equal(b, user.PhoneNumber, resUser.PhoneNumber)
		assert.Equal(b, user.Email, resUser.Email)
		assert.Equal(b, user.Address, resUser.Address)
		assert.Equal(b, user.Company, resUser.Company)
	}
}

func BenchmarkNova_HandleCreateUserParallel(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateUserParallel
	// Test Purpose: Benchmark HandleCreateUser create user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			/* create userId */
			// request content
			url := server.URL + "/nova/v1/user/userId"
			// request create userId
			wUserId := httptest.NewRecorder()
			reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wUserId, reqUserId)
			// return response
			var resUserId UserID
			err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusOK, wUserId.Code)
			assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
			assert.NoError(b, uuid.Validate(resUserId.UserId))
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
				b.Errorf("error marshal user: %v", err)
			}
			// request create user
			wUser := httptest.NewRecorder()
			reqUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqUser.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wUser, reqUser)
			// return response
			var resUser User
			err = json.Unmarshal(wUser.Body.Bytes(), &resUser)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wUser.Code)
			assert.Equal(b, "application/json", wUser.Header().Get("Content-Type"))
			assert.Equal(b, user.UserId, resUser.UserId)
			assert.Equal(b, user.Username, resUser.Username)
			assert.Equal(b, user.Password, resUser.Password)
			assert.Equal(b, user.PhoneNumber, resUser.PhoneNumber)
			assert.Equal(b, user.Email, resUser.Email)
			assert.Equal(b, user.Address, resUser.Address)
			assert.Equal(b, user.Company, resUser.Company)
		}
	})
}

func TestNova_HandleDeleteUser(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestNova_HandleDeleteUser
	// Test Purpose: Test HandleDeleteUser delete user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send DeleteUser request with userId by using DELETE method
	// 6. receive DeleteUser request by using 204 No Content Code
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
	wCreateUser := httptest.NewRecorder()
	reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqCreateUser.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wCreateUser, reqCreateUser)
	// return response
	var resUser User
	err = json.Unmarshal(wCreateUser.Body.Bytes(), &resUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wCreateUser.Code)
	assert.Equal(t, "application/json", wCreateUser.Header().Get("Content-Type"))
	assert.Equal(t, user.UserId, resUser.UserId)
	assert.Equal(t, user.Username, resUser.Username)
	assert.Equal(t, user.Password, resUser.Password)
	assert.Equal(t, user.PhoneNumber, resUser.PhoneNumber)
	assert.Equal(t, user.Email, resUser.Email)
	assert.Equal(t, user.Address, resUser.Address)
	assert.Equal(t, user.Company, resUser.Company)
	/* delete user */
	// request content
	url = server.URL + "/nova/v1/user"
	// request delete user
	wDeleteUser := httptest.NewRecorder()
	reqDeleteUser, err := http.NewRequest(http.MethodDelete, url+"/"+resUserId.UserId, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wDeleteUser, reqDeleteUser)
	// validate response
	assert.Equal(t, http.StatusNoContent, wDeleteUser.Code)
}

func BenchmarkNova_HandleDeleteUser(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteUser
	// Test Purpose: Benchmark HandleDeleteUser delete user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send DeleteUser request with userId by using DELETE method
	// 6. receive DeleteUser request by using 204 No Content Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		/* create userId */
		// request content
		url := server.URL + "/nova/v1/user/userId"
		// request create userId
		wUserId := httptest.NewRecorder()
		reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wUserId, reqUserId)
		// return response
		var resUserId UserID
		err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusOK, wUserId.Code)
		assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
		assert.NoError(b, uuid.Validate(resUserId.UserId))
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
			b.Errorf("error marshal user: %v", err)
		}
		// request create user
		wCreateUser := httptest.NewRecorder()
		reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqCreateUser.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wCreateUser, reqCreateUser)
		// return response
		var resUser User
		err = json.Unmarshal(wCreateUser.Body.Bytes(), &resUser)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wCreateUser.Code)
		assert.Equal(b, "application/json", wCreateUser.Header().Get("Content-Type"))
		assert.Equal(b, user.UserId, resUser.UserId)
		assert.Equal(b, user.Username, resUser.Username)
		assert.Equal(b, user.Password, resUser.Password)
		assert.Equal(b, user.PhoneNumber, resUser.PhoneNumber)
		assert.Equal(b, user.Email, resUser.Email)
		assert.Equal(b, user.Address, resUser.Address)
		assert.Equal(b, user.Company, resUser.Company)
		/* delete user */
		// request content
		url = server.URL + "/nova/v1/user"
		// request delete user
		wDeleteUser := httptest.NewRecorder()
		reqDeleteUser, err := http.NewRequest(http.MethodDelete, url+"/"+resUserId.UserId, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wDeleteUser, reqDeleteUser)
		// validate response
		assert.Equal(b, http.StatusNoContent, wDeleteUser.Code)
	}
}

func BenchmarkNova_HandleDeleteUserParallel(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteUserParallel
	// Test Purpose: Benchmark HandleDeleteUser delete user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send DeleteUser request with userId by using DELETE method
	// 6. receive DeleteUser request by using 204 No Content Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			/* create userId */
			// request content
			url := server.URL + "/nova/v1/user/userId"
			// request create userId
			wUserId := httptest.NewRecorder()
			reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wUserId, reqUserId)
			// return response
			var resUserId UserID
			err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusOK, wUserId.Code)
			assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
			assert.NoError(b, uuid.Validate(resUserId.UserId))
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
				b.Errorf("error marshal user: %v", err)
			}
			// request create user
			wCreateUser := httptest.NewRecorder()
			reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqCreateUser.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wCreateUser, reqCreateUser)
			// return response
			var resUser User
			err = json.Unmarshal(wCreateUser.Body.Bytes(), &resUser)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wCreateUser.Code)
			assert.Equal(b, "application/json", wCreateUser.Header().Get("Content-Type"))
			assert.Equal(b, user.UserId, resUser.UserId)
			assert.Equal(b, user.Username, resUser.Username)
			assert.Equal(b, user.Password, resUser.Password)
			assert.Equal(b, user.PhoneNumber, resUser.PhoneNumber)
			assert.Equal(b, user.Email, resUser.Email)
			assert.Equal(b, user.Address, resUser.Address)
			assert.Equal(b, user.Company, resUser.Company)
			/* delete user */
			// request content
			url = server.URL + "/nova/v1/user"
			// request delete user
			wDeleteUser := httptest.NewRecorder()
			reqDeleteUser, err := http.NewRequest(http.MethodDelete, url+"/"+resUserId.UserId, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wDeleteUser, reqDeleteUser)
			// validate response
			assert.Equal(b, http.StatusNoContent, wDeleteUser.Code)
		}
	})
}

func TestNova_HandleQueryUserUser(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestNova_HandleQueryUserUser
	// Test Purpose: Test HandleQueryUserUser query user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send QueryUser request with userId by using GET method
	// 6. receive QueryUser request by using 200 OK Code
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
	wCreateUser := httptest.NewRecorder()
	reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqCreateUser.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wCreateUser, reqCreateUser)
	// return response
	var resCreateUser User
	err = json.Unmarshal(wCreateUser.Body.Bytes(), &resCreateUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wCreateUser.Code)
	assert.Equal(t, "application/json", wCreateUser.Header().Get("Content-Type"))
	assert.Equal(t, user.UserId, resCreateUser.UserId)
	assert.Equal(t, user.Username, resCreateUser.Username)
	assert.Equal(t, user.Password, resCreateUser.Password)
	assert.Equal(t, user.PhoneNumber, resCreateUser.PhoneNumber)
	assert.Equal(t, user.Email, resCreateUser.Email)
	assert.Equal(t, user.Address, resCreateUser.Address)
	assert.Equal(t, user.Company, resCreateUser.Company)
	/* query user */
	// request content
	url = server.URL + "/nova/v1/user"
	// request query user
	wQueryUser := httptest.NewRecorder()
	reqQueryUser, err := http.NewRequest(http.MethodGet, url+"/"+resUserId.UserId, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wQueryUser, reqQueryUser)
	// return response
	var resQueryUser User
	err = json.Unmarshal(wQueryUser.Body.Bytes(), &resQueryUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusOK, wQueryUser.Code)
	assert.Equal(t, "application/json", wQueryUser.Header().Get("Content-Type"))
	assert.Equal(t, user.UserId, resQueryUser.UserId)
	assert.Equal(t, user.Username, resQueryUser.Username)
	assert.Equal(t, user.Password, resQueryUser.Password)
	assert.Equal(t, user.PhoneNumber, resQueryUser.PhoneNumber)
	assert.Equal(t, user.Email, resQueryUser.Email)
	assert.Equal(t, user.Address, resQueryUser.Address)
	assert.Equal(t, user.Company, resQueryUser.Company)
}

func BenchmarkNova_HandleQueryUser(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleQueryUser
	// Test Purpose: Benchmark HandleQueryUserUser query user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send QueryUser request with userId by using GET method
	// 6. receive QueryUser request by using 200 OK Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		/* create userId */
		// request content
		url := server.URL + "/nova/v1/user/userId"
		// request create userId
		wUserId := httptest.NewRecorder()
		reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wUserId, reqUserId)
		// return response
		var resUserId UserID
		err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusOK, wUserId.Code)
		assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
		assert.NoError(b, uuid.Validate(resUserId.UserId))
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
			b.Errorf("error marshal user: %v", err)
		}
		// request create user
		wCreateUser := httptest.NewRecorder()
		reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqCreateUser.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wCreateUser, reqCreateUser)
		// return response
		var resCreateUser User
		err = json.Unmarshal(wCreateUser.Body.Bytes(), &resCreateUser)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wCreateUser.Code)
		assert.Equal(b, "application/json", wCreateUser.Header().Get("Content-Type"))
		assert.Equal(b, user.UserId, resCreateUser.UserId)
		assert.Equal(b, user.Username, resCreateUser.Username)
		assert.Equal(b, user.Password, resCreateUser.Password)
		assert.Equal(b, user.PhoneNumber, resCreateUser.PhoneNumber)
		assert.Equal(b, user.Email, resCreateUser.Email)
		assert.Equal(b, user.Address, resCreateUser.Address)
		assert.Equal(b, user.Company, resCreateUser.Company)
		/* query user */
		// request content
		url = server.URL + "/nova/v1/user"
		// request query user
		wQueryUser := httptest.NewRecorder()
		reqQueryUser, err := http.NewRequest(http.MethodGet, url+"/"+resUserId.UserId, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wQueryUser, reqQueryUser)
		// return response
		var resQueryUser User
		err = json.Unmarshal(wQueryUser.Body.Bytes(), &resQueryUser)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusOK, wQueryUser.Code)
		assert.Equal(b, "application/json", wQueryUser.Header().Get("Content-Type"))
		assert.Equal(b, user.UserId, resQueryUser.UserId)
		assert.Equal(b, user.Username, resQueryUser.Username)
		assert.Equal(b, user.Password, resQueryUser.Password)
		assert.Equal(b, user.PhoneNumber, resQueryUser.PhoneNumber)
		assert.Equal(b, user.Email, resQueryUser.Email)
		assert.Equal(b, user.Address, resQueryUser.Address)
		assert.Equal(b, user.Company, resQueryUser.Company)
	}
}

func BenchmarkNova_HandleQueryUserParallel(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleQueryUserParallel
	// Test Purpose: Benchmark HandleQueryUserUser query user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send QueryUser request with userId by using GET method
	// 6. receive QueryUser request by using 200 OK Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			/* create userId */
			// request content
			url := server.URL + "/nova/v1/user/userId"
			// request create userId
			wUserId := httptest.NewRecorder()
			reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wUserId, reqUserId)
			// return response
			var resUserId UserID
			err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusOK, wUserId.Code)
			assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
			assert.NoError(b, uuid.Validate(resUserId.UserId))
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
				b.Errorf("error marshal user: %v", err)
			}
			// request create user
			wCreateUser := httptest.NewRecorder()
			reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqCreateUser.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wCreateUser, reqCreateUser)
			// return response
			var resCreateUser User
			err = json.Unmarshal(wCreateUser.Body.Bytes(), &resCreateUser)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wCreateUser.Code)
			assert.Equal(b, "application/json", wCreateUser.Header().Get("Content-Type"))
			assert.Equal(b, user.UserId, resCreateUser.UserId)
			assert.Equal(b, user.Username, resCreateUser.Username)
			assert.Equal(b, user.Password, resCreateUser.Password)
			assert.Equal(b, user.PhoneNumber, resCreateUser.PhoneNumber)
			assert.Equal(b, user.Email, resCreateUser.Email)
			assert.Equal(b, user.Address, resCreateUser.Address)
			assert.Equal(b, user.Company, resCreateUser.Company)
			/* query user */
			// request content
			url = server.URL + "/nova/v1/user"
			// request query user
			wQueryUser := httptest.NewRecorder()
			reqQueryUser, err := http.NewRequest(http.MethodGet, url+"/"+resUserId.UserId, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wQueryUser, reqQueryUser)
			// return response
			var resQueryUser User
			err = json.Unmarshal(wQueryUser.Body.Bytes(), &resQueryUser)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusOK, wQueryUser.Code)
			assert.Equal(b, "application/json", wQueryUser.Header().Get("Content-Type"))
			assert.Equal(b, user.UserId, resQueryUser.UserId)
			assert.Equal(b, user.Username, resQueryUser.Username)
			assert.Equal(b, user.Password, resQueryUser.Password)
			assert.Equal(b, user.PhoneNumber, resQueryUser.PhoneNumber)
			assert.Equal(b, user.Email, resQueryUser.Email)
			assert.Equal(b, user.Address, resQueryUser.Address)
			assert.Equal(b, user.Company, resQueryUser.Company)
		}
	})
}

func TestNova_HandleUpdateUser(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestNova_HandleUpdateUser
	// Test Purpose: Test HandleUpdateUser update user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send UpdateUser request with userId by using PUT method
	// 6. receive UpdateUser request by using 200 OK Code
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
	wCreateUser := httptest.NewRecorder()
	reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqCreateUser.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wCreateUser, reqCreateUser)
	// return response
	var resCreateUser User
	err = json.Unmarshal(wCreateUser.Body.Bytes(), &resCreateUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wCreateUser.Code)
	assert.Equal(t, "application/json", wCreateUser.Header().Get("Content-Type"))
	assert.Equal(t, user.UserId, resCreateUser.UserId)
	assert.Equal(t, user.Username, resCreateUser.Username)
	assert.Equal(t, user.Password, resCreateUser.Password)
	assert.Equal(t, user.PhoneNumber, resCreateUser.PhoneNumber)
	assert.Equal(t, user.Email, resCreateUser.Email)
	assert.Equal(t, user.Address, resCreateUser.Address)
	assert.Equal(t, user.Company, resCreateUser.Company)
	/* update user */
	// request content
	url = server.URL + "/nova/v1/user"
	userNew := User{
		UserId:      resUserId.UserId,
		Username:    "bob",
		Password:    "888888",
		PhoneNumber: "+2839822",
		Email:       "bob@gmail.com",
		Address:     "No.101, New Street, Los Angle, USA",
		Company:     "Microsoft",
	}
	bodyNew, err := json.Marshal(userNew)
	if err != nil {
		t.Errorf("error marshal user: %v", err)
	}
	// request update user
	wUpdateUser := httptest.NewRecorder()
	reqUpdateUser, err := http.NewRequest(http.MethodPut, url+"/"+resUserId.UserId, bytes.NewReader(bodyNew))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wUpdateUser, reqUpdateUser)
	// return response
	var resUpdateUser User
	err = json.Unmarshal(wUpdateUser.Body.Bytes(), &resUpdateUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusOK, wUpdateUser.Code)
	assert.Equal(t, "application/json", wUpdateUser.Header().Get("Content-Type"))
	assert.Equal(t, userNew.UserId, resUpdateUser.UserId)
	assert.Equal(t, userNew.Username, resUpdateUser.Username)
	assert.Equal(t, userNew.Password, resUpdateUser.Password)
	assert.Equal(t, userNew.PhoneNumber, resUpdateUser.PhoneNumber)
	assert.Equal(t, userNew.Email, resUpdateUser.Email)
	assert.Equal(t, userNew.Address, resUpdateUser.Address)
	assert.Equal(t, userNew.Company, resUpdateUser.Company)
}

func BenchmarkNova_HandleUpdateUser(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleUpdateUser
	// Test Purpose: Benchmark HandleUpdateUser update user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send UpdateUser request with userId by using PUT method
	// 6. receive UpdateUser request by using 200 OK Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		/* create userId */
		// request content
		url := server.URL + "/nova/v1/user/userId"
		// request create userId
		wUserId := httptest.NewRecorder()
		reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wUserId, reqUserId)
		// return response
		var resUserId UserID
		err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusOK, wUserId.Code)
		assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
		assert.NoError(b, uuid.Validate(resUserId.UserId))
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
			b.Errorf("error marshal user: %v", err)
		}
		// request create user
		wCreateUser := httptest.NewRecorder()
		reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqCreateUser.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wCreateUser, reqCreateUser)
		// return response
		var resCreateUser User
		err = json.Unmarshal(wCreateUser.Body.Bytes(), &resCreateUser)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wCreateUser.Code)
		assert.Equal(b, "application/json", wCreateUser.Header().Get("Content-Type"))
		assert.Equal(b, user.UserId, resCreateUser.UserId)
		assert.Equal(b, user.Username, resCreateUser.Username)
		assert.Equal(b, user.Password, resCreateUser.Password)
		assert.Equal(b, user.PhoneNumber, resCreateUser.PhoneNumber)
		assert.Equal(b, user.Email, resCreateUser.Email)
		assert.Equal(b, user.Address, resCreateUser.Address)
		assert.Equal(b, user.Company, resCreateUser.Company)
		/* update user */
		// request content
		url = server.URL + "/nova/v1/user"
		userNew := User{
			UserId:      resUserId.UserId,
			Username:    "bob",
			Password:    "888888",
			PhoneNumber: "+2839822",
			Email:       "bob@gmail.com",
			Address:     "No.101, New Street, Los Angle, USA",
			Company:     "Microsoft",
		}
		bodyNew, err := json.Marshal(userNew)
		if err != nil {
			b.Errorf("error marshal user: %v", err)
		}
		// request update user
		wUpdateUser := httptest.NewRecorder()
		reqUpdateUser, err := http.NewRequest(http.MethodPut, url+"/"+resUserId.UserId, bytes.NewReader(bodyNew))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wUpdateUser, reqUpdateUser)
		// return response
		var resUpdateUser User
		err = json.Unmarshal(wUpdateUser.Body.Bytes(), &resUpdateUser)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusOK, wUpdateUser.Code)
		assert.Equal(b, "application/json", wUpdateUser.Header().Get("Content-Type"))
		assert.Equal(b, userNew.UserId, resUpdateUser.UserId)
		assert.Equal(b, userNew.Username, resUpdateUser.Username)
		assert.Equal(b, userNew.Password, resUpdateUser.Password)
		assert.Equal(b, userNew.PhoneNumber, resUpdateUser.PhoneNumber)
		assert.Equal(b, userNew.Email, resUpdateUser.Email)
		assert.Equal(b, userNew.Address, resUpdateUser.Address)
		assert.Equal(b, userNew.Company, resUpdateUser.Company)
	}
}

func BenchmarkNova_HandleUpdateUserParallel(b *testing.B) {
	/*--------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleUpdateUserParallel
	// Test Purpose: Benchmark HandleUpdateUser update user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send UpdateUser request with userId by using PUT method
	// 6. receive UpdateUser request by using 200 OK Code
	----------------------------------------------------------------------------------*/
	// start http test service
	server, router := startTestService()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			/* create userId */
			// request content
			url := server.URL + "/nova/v1/user/userId"
			// request create userId
			wUserId := httptest.NewRecorder()
			reqUserId, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wUserId, reqUserId)
			// return response
			var resUserId UserID
			err = json.Unmarshal(wUserId.Body.Bytes(), &resUserId)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusOK, wUserId.Code)
			assert.Equal(b, "application/json", wUserId.Header().Get("Content-Type"))
			assert.NoError(b, uuid.Validate(resUserId.UserId))
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
				b.Errorf("error marshal user: %v", err)
			}
			// request create user
			wCreateUser := httptest.NewRecorder()
			reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqCreateUser.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wCreateUser, reqCreateUser)
			// return response
			var resCreateUser User
			err = json.Unmarshal(wCreateUser.Body.Bytes(), &resCreateUser)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wCreateUser.Code)
			assert.Equal(b, "application/json", wCreateUser.Header().Get("Content-Type"))
			assert.Equal(b, user.UserId, resCreateUser.UserId)
			assert.Equal(b, user.Username, resCreateUser.Username)
			assert.Equal(b, user.Password, resCreateUser.Password)
			assert.Equal(b, user.PhoneNumber, resCreateUser.PhoneNumber)
			assert.Equal(b, user.Email, resCreateUser.Email)
			assert.Equal(b, user.Address, resCreateUser.Address)
			assert.Equal(b, user.Company, resCreateUser.Company)
			/* update user */
			// request content
			url = server.URL + "/nova/v1/user"
			userNew := User{
				UserId:      resUserId.UserId,
				Username:    "bob",
				Password:    "888888",
				PhoneNumber: "+2839822",
				Email:       "bob@gmail.com",
				Address:     "No.101, New Street, Los Angle, USA",
				Company:     "Microsoft",
			}
			bodyNew, err := json.Marshal(userNew)
			if err != nil {
				b.Errorf("error marshal user: %v", err)
			}
			// request update user
			wUpdateUser := httptest.NewRecorder()
			reqUpdateUser, err := http.NewRequest(http.MethodPut, url+"/"+resUserId.UserId, bytes.NewReader(bodyNew))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wUpdateUser, reqUpdateUser)
			// return response
			var resUpdateUser User
			err = json.Unmarshal(wUpdateUser.Body.Bytes(), &resUpdateUser)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusOK, wUpdateUser.Code)
			assert.Equal(b, "application/json", wUpdateUser.Header().Get("Content-Type"))
			assert.Equal(b, userNew.UserId, resUpdateUser.UserId)
			assert.Equal(b, userNew.Username, resUpdateUser.Username)
			assert.Equal(b, userNew.Password, resUpdateUser.Password)
			assert.Equal(b, userNew.PhoneNumber, resUpdateUser.PhoneNumber)
			assert.Equal(b, userNew.Email, resUpdateUser.Email)
			assert.Equal(b, userNew.Address, resUpdateUser.Address)
			assert.Equal(b, userNew.Company, resUpdateUser.Company)
		}
	})
}

func TestNova_HandleModifyUser(t *testing.T) {
	/*--------------------------------------------------------------------------------
	// Test Case: TestNova_HandleModifyUser
	// Test Purpose: Test HandleModifyUser modify user
	// Test Steps:
	// 1. send CreateUserId request by using POST method
	// 2. receive CreateUserId response with created userId by using 200 OK Code
	// 3. send CreateUser request with user information by using POST method
	// 4. receive CreateUser response with user information by using 201 Created Code
	// 5. send ModifyUser request with userId by using PATCH method
	// 6. receive ModifyUser request by using 200 OK Code
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
	wCreateUser := httptest.NewRecorder()
	reqCreateUser, err := http.NewRequest(http.MethodPost, url+"/"+resUserId.UserId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqCreateUser.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wCreateUser, reqCreateUser)
	// return response
	var resCreateUser User
	err = json.Unmarshal(wCreateUser.Body.Bytes(), &resCreateUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wCreateUser.Code)
	assert.Equal(t, "application/json", wCreateUser.Header().Get("Content-Type"))
	assert.Equal(t, user.UserId, resCreateUser.UserId)
	assert.Equal(t, user.Username, resCreateUser.Username)
	assert.Equal(t, user.Password, resCreateUser.Password)
	assert.Equal(t, user.PhoneNumber, resCreateUser.PhoneNumber)
	assert.Equal(t, user.Email, resCreateUser.Email)
	assert.Equal(t, user.Address, resCreateUser.Address)
	assert.Equal(t, user.Company, resCreateUser.Company)
	/* modify user */
	// request content
	url = server.URL + "/nova/v1/user"
	userNew := User{
		UserId:      resUserId.UserId,
		Username:    "alice",
		Password:    "123456",
		PhoneNumber: "+1412387",
		Company:     "Microsoft",
	}
	bodyNew, err := json.Marshal(userNew)
	if err != nil {
		t.Errorf("error marshal user: %v", err)
	}
	// request modify user
	wModifyUser := httptest.NewRecorder()
	reqModifyUser, err := http.NewRequest(http.MethodPatch, url+"/"+resUserId.UserId, bytes.NewReader(bodyNew))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wModifyUser, reqModifyUser)
	// return response
	var resModifyUser User
	err = json.Unmarshal(wModifyUser.Body.Bytes(), &resModifyUser)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusOK, wModifyUser.Code)
	assert.Equal(t, "application/json", wModifyUser.Header().Get("Content-Type"))
	assert.Equal(t, userNew.UserId, resModifyUser.UserId)
	assert.Equal(t, userNew.Username, resModifyUser.Username)
	assert.Equal(t, userNew.Password, resModifyUser.Password)
	assert.Equal(t, userNew.PhoneNumber, resModifyUser.PhoneNumber)
	assert.Equal(t, user.Email, resModifyUser.Email)
	assert.Equal(t, user.Address, resModifyUser.Address)
	assert.Equal(t, userNew.Company, resModifyUser.Company)
}

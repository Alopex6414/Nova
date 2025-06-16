package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupQuestionTestRouter() *gin.Engine {
	// create Nova instance
	nova := New()
	// initialize Nova instance
	nova.Init()
	// apply default Gin service
	router := gin.Default()
	// apply Gin logger & recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// create router group for nova
	novaService := router.Group("nova/v1")
	{
		novaService.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "hello Nova\n") })
		/* question management */
		// questionId related
		novaService.POST("/question/Id", nova.HandleCreateQuestionId)
		// question related
		novaService.POST("/question/:Id", nova.HandleCreateQuestion)
		novaService.PUT("/question/:Id", nova.HandleUpdateQuestion)
		novaService.DELETE("/question/:Id", nova.HandleDeleteQuestion)
		novaService.PATCH("/question/:Id", nova.HandleModifyQuestion)
		novaService.GET("/question/:Id", nova.HandleQueryQuestion)
	}
	return router
}

func startQuestionTestService() (*httptest.Server, *gin.Engine) {
	router := setupQuestionTestRouter()
	return httptest.NewServer(router), router
}

func resetQuestionTestCase() error {
	// remove database
	err := os.Remove("nova.db")
	if err != nil {
		return err
	}
	// remove logs
	err = os.RemoveAll("logs/")
	if err != nil {
		return err
	}
	return nil
}

func TestNova_HandleCreateQuestionId(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateQuestionId
	// Test Purpose: Test HandleCreateQuestionId create questionId
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	-----------------------------------------------------------------------------------------*/
	// reset test case
	_ = resetQuestionTestCase()
	// start http test service
	server, router := startQuestionTestService()
	defer server.Close()
	// request content
	url := server.URL + "/nova/v1/question/Id"
	// request create questionId
	w := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(w, request)
	// return response
	var response string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.NoError(t, uuid.Validate(response))
}

package app

import (
	"bytes"
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

func BenchmarkNova_HandleCreateQuestionId(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestionId
	// Test Purpose: Benchmark HandleCreateQuestionId create questionId
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	-----------------------------------------------------------------------------------------*/
	// reset test case
	_ = resetQuestionTestCase()
	// start http test service
	server, router := startQuestionTestService()
	defer server.Close()
	// start benchmark test
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// request content
		url := server.URL + "/nova/v1/question/Id"
		// request create questionId
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(w, request)
		// return response
		var response string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.NoError(b, uuid.Validate(response))
	}
}

func BenchmarkNova_HandleCreateQuestionIdParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestionId (Parallel)
	// Test Purpose: Benchmark HandleCreateQuestionId create questionId
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	-----------------------------------------------------------------------------------------*/
	// reset test case
	_ = resetQuestionTestCase()
	// start http test service
	server, router := startQuestionTestService()
	defer server.Close()
	// start benchmark test
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// request content
			url := server.URL + "/nova/v1/question/Id"
			// request create questionId
			w := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(w, request)
			// return response
			var response string
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, w.Code)
			assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
			assert.NoError(b, uuid.Validate(response))
		}
	})
}

func TestNova_HandleCreateQuestion(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateQuestion
	// Test Purpose: Test HandleCreateQuestion create question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	-----------------------------------------------------------------------------------------*/
	// reset test case
	_ = resetQuestionTestCase()
	// start http test service
	server, router := startQuestionTestService()
	defer server.Close()
	/* create questionId */
	// request content
	url := server.URL + "/nova/v1/question/Id"
	// request create questionId
	wQuestionId := httptest.NewRecorder()
	reqQuestionId, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wQuestionId, reqQuestionId)
	// return response
	var reQuestionId string
	err = json.Unmarshal(wQuestionId.Body.Bytes(), &reQuestionId)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wQuestionId.Code)
	assert.Equal(t, "application/json", wQuestionId.Header().Get("Content-Type"))
	assert.NoError(t, uuid.Validate(reQuestionId))
	/* create question */
	url = server.URL + "/nova/v1/question"
	question := QuestionSingleChoice{
		Id:    reQuestionId,
		Title: "What's the sweetest fruit?",
		Answers: []QuestionAnswer{
			QuestionAnswer{
				"A",
				"apple",
			},
			QuestionAnswer{
				"B",
				"watermelon",
			},
			QuestionAnswer{
				"C",
				"orange",
			},
			QuestionAnswer{
				"D",
				"peach",
			},
		},
		StandardAnswer: QuestionAnswer{
			"B",
			"watermelon",
		},
	}
	body, err := json.Marshal(question)
	if err != nil {
		t.Errorf("error marshal question: %v", err)
	}
	// request create user
	wQuestion := httptest.NewRecorder()
	reqQuestion, err := http.NewRequest(http.MethodPost, url+"/"+reQuestionId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqQuestion.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wQuestion, reqQuestion)
	// return response
	var resQuestion QuestionSingleChoice
	err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wQuestion.Code)
	assert.Equal(t, "application/json", wQuestion.Header().Get("Content-Type"))
	assert.Equal(t, question.Id, resQuestion.Id)
	assert.Equal(t, question.Title, resQuestion.Title)
	assert.Equal(t, question.Answers, resQuestion.Answers)
	assert.Equal(t, question.StandardAnswer, resQuestion.StandardAnswer)
}

func BenchmarkNova_HandleCreateQuestion(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion
	// Test Purpose: Benchmark HandleCreateQuestion create question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	-----------------------------------------------------------------------------------------*/
	// reset test case
	_ = resetQuestionTestCase()
	// start http test service
	server, router := startQuestionTestService()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		/* create questionId */
		// request content
		url := server.URL + "/nova/v1/question/Id"
		// request create questionId
		wQuestionId := httptest.NewRecorder()
		reqQuestionId, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wQuestionId, reqQuestionId)
		// return response
		var reQuestionId string
		err = json.Unmarshal(wQuestionId.Body.Bytes(), &reQuestionId)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wQuestionId.Code)
		assert.Equal(b, "application/json", wQuestionId.Header().Get("Content-Type"))
		assert.NoError(b, uuid.Validate(reQuestionId))
		/* create question */
		url = server.URL + "/nova/v1/question"
		question := QuestionSingleChoice{
			Id:    reQuestionId,
			Title: "What's the sweetest fruit?",
			Answers: []QuestionAnswer{
				QuestionAnswer{
					"A",
					"apple",
				},
				QuestionAnswer{
					"B",
					"watermelon",
				},
				QuestionAnswer{
					"C",
					"orange",
				},
				QuestionAnswer{
					"D",
					"peach",
				},
			},
			StandardAnswer: QuestionAnswer{
				"B",
				"watermelon",
			},
		}
		body, err := json.Marshal(question)
		if err != nil {
			b.Errorf("error marshal question: %v", err)
		}
		// request create user
		wQuestion := httptest.NewRecorder()
		reqQuestion, err := http.NewRequest(http.MethodPost, url+"/"+reQuestionId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqQuestion.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wQuestion, reqQuestion)
		// return response
		var resQuestion QuestionSingleChoice
		err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wQuestion.Code)
		assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
		assert.Equal(b, question.Id, resQuestion.Id)
		assert.Equal(b, question.Title, resQuestion.Title)
		assert.Equal(b, question.Answers, resQuestion.Answers)
		assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
	}
}

func BenchmarkNova_HandleCreateQuestionParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (Parallel)
	// Test Purpose: Benchmark HandleCreateQuestion create question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	-----------------------------------------------------------------------------------------*/
	// reset test case
	_ = resetQuestionTestCase()
	// start http test service
	server, router := startQuestionTestService()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			/* create questionId */
			// request content
			url := server.URL + "/nova/v1/question/Id"
			// request create questionId
			wQuestionId := httptest.NewRecorder()
			reqQuestionId, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wQuestionId, reqQuestionId)
			// return response
			var reQuestionId string
			err = json.Unmarshal(wQuestionId.Body.Bytes(), &reQuestionId)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wQuestionId.Code)
			assert.Equal(b, "application/json", wQuestionId.Header().Get("Content-Type"))
			assert.NoError(b, uuid.Validate(reQuestionId))
			/* create question */
			url = server.URL + "/nova/v1/question"
			question := QuestionSingleChoice{
				Id:    reQuestionId,
				Title: "What's the sweetest fruit?",
				Answers: []QuestionAnswer{
					QuestionAnswer{
						"A",
						"apple",
					},
					QuestionAnswer{
						"B",
						"watermelon",
					},
					QuestionAnswer{
						"C",
						"orange",
					},
					QuestionAnswer{
						"D",
						"peach",
					},
				},
				StandardAnswer: QuestionAnswer{
					"B",
					"watermelon",
				},
			}
			body, err := json.Marshal(question)
			if err != nil {
				b.Errorf("error marshal question: %v", err)
			}
			// request create user
			wQuestion := httptest.NewRecorder()
			reqQuestion, err := http.NewRequest(http.MethodPost, url+"/"+reQuestionId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqQuestion.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wQuestion, reqQuestion)
			// return response
			var resQuestion QuestionSingleChoice
			err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wQuestion.Code)
			assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
			assert.Equal(b, question.Id, resQuestion.Id)
			assert.Equal(b, question.Title, resQuestion.Title)
			assert.Equal(b, question.Answers, resQuestion.Answers)
			assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
		}
	})
}

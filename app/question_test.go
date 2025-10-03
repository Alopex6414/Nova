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
		novaService.POST("/question/single-choice/:Id", nova.HandleCreateQuestionSingleChoice)
		novaService.PUT("/question/single-choice/:Id", nova.HandleUpdateQuestionSingleChoice)
		novaService.DELETE("/question/single-choice/:Id", nova.HandleDeleteQuestionSingleChoice)
		novaService.PATCH("/question/single-choice/:Id", nova.HandleModifyQuestionSingleChoice)
		novaService.GET("/question/single-choice/:Id", nova.HandleQueryQuestionSingleChoice)
		novaService.POST("/question/multiple-choice/:Id", nova.HandleCreateQuestionMultipleChoice)
		novaService.PUT("/question/multiple-choice/:Id", nova.HandleUpdateQuestionMultipleChoice)
		novaService.DELETE("/question/multiple-choice/:Id", nova.HandleDeleteQuestionMultipleChoice)
		novaService.PATCH("/question/multiple-choice/:Id", nova.HandleModifyQuestionMultipleChoice)
		novaService.GET("/question/multiple-choice/:Id", nova.HandleQueryQuestionMultipleChoice)
		novaService.POST("/question/judgement/:Id", nova.HandleCreateQuestionJudgement)
		novaService.PUT("/question/judgement/:Id", nova.HandleUpdateQuestionJudgement)
		novaService.DELETE("/question/judgement/:Id", nova.HandleDeleteQuestionJudgement)
		novaService.PATCH("/question/judgement/:Id", nova.HandleModifyQuestionJudgement)
		novaService.GET("/question/judgement/:Id", nova.HandleQueryQuestionJudgement)
		novaService.POST("/question/essay/:Id", nova.HandleCreateQuestionEssay)
		novaService.PUT("/question/essay/:Id", nova.HandleUpdateQuestionEssay)
		novaService.DELETE("/question/essay/:Id", nova.HandleDeleteQuestionEssay)
		novaService.PATCH("/question/essay/:Id", nova.HandleModifyQuestionEssay)
		novaService.GET("/question/essay/:Id", nova.HandleQueryQuestionEssay)
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

func TestNova_HandleCreateQuestionSingleChoice(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateQuestion (single-choice)
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
	url = server.URL + "/nova/v1/question/single-choice"
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

func BenchmarkNova_HandleCreateQuestionSingleChoice(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (single-choice)
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
		url = server.URL + "/nova/v1/question/single-choice"
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

func BenchmarkNova_HandleCreateQuestionSingleChoiceParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (single-choice) (Parallel)
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
			url = server.URL + "/nova/v1/question/single-choice"
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

func TestNova_HandleCreateQuestionMultipleChoice(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateQuestion (multiple-choice)
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
	url = server.URL + "/nova/v1/question/multiple-choice"
	question := QuestionMultipleChoice{
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
		StandardAnswers: []QuestionAnswer{
			{
				"B",
				"watermelon",
			},
			{
				"D",
				"peach",
			},
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
	var resQuestion QuestionMultipleChoice
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
	assert.Equal(t, question.StandardAnswers, resQuestion.StandardAnswers)
}

func BenchmarkNova_HandleCreateQuestionMultipleChoice(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (multiple-choice)
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
		url = server.URL + "/nova/v1/question/multiple-choice"
		question := QuestionMultipleChoice{
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
			StandardAnswers: []QuestionAnswer{
				{
					"B",
					"watermelon",
				},
				{
					"D",
					"peach",
				},
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
		var resQuestion QuestionMultipleChoice
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
		assert.Equal(b, question.StandardAnswers, resQuestion.StandardAnswers)
	}
}

func BenchmarkNova_HandleCreateQuestionMultipleChoiceParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (multiple-choice) (Parallel)
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
			url = server.URL + "/nova/v1/question/multiple-choice"
			question := QuestionMultipleChoice{
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
				StandardAnswers: []QuestionAnswer{
					{
						"B",
						"watermelon",
					},
					{
						"D",
						"peach",
					},
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
			var resQuestion QuestionMultipleChoice
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
			assert.Equal(b, question.StandardAnswers, resQuestion.StandardAnswers)
		}
	})
}

func TestNova_HandleCreateQuestionJudgement(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateQuestion (judgement)
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
	url = server.URL + "/nova/v1/question/judgement"
	question := QuestionJudgement{
		Id:             reQuestionId,
		Title:          "What's the sweetest fruit?",
		Answer:         true,
		StandardAnswer: false,
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
	var resQuestion QuestionJudgement
	err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wQuestion.Code)
	assert.Equal(t, "application/json", wQuestion.Header().Get("Content-Type"))
	assert.Equal(t, question.Id, resQuestion.Id)
	assert.Equal(t, question.Title, resQuestion.Title)
	assert.Equal(t, question.Answer, resQuestion.Answer)
	assert.Equal(t, question.StandardAnswer, resQuestion.StandardAnswer)
}

func BenchmarkNova_HandleCreateQuestionJudgement(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (judgement)
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
		url = server.URL + "/nova/v1/question/judgement"
		question := QuestionJudgement{
			Id:             reQuestionId,
			Title:          "What's the sweetest fruit?",
			Answer:         true,
			StandardAnswer: false,
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
		var resQuestion QuestionJudgement
		err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wQuestion.Code)
		assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
		assert.Equal(b, question.Id, resQuestion.Id)
		assert.Equal(b, question.Title, resQuestion.Title)
		assert.Equal(b, question.Answer, resQuestion.Answer)
		assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
	}
}

func BenchmarkNova_HandleCreateQuestionJudgementParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (judgement)
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
			url = server.URL + "/nova/v1/question/judgement"
			question := QuestionJudgement{
				Id:             reQuestionId,
				Title:          "What's the sweetest fruit?",
				Answer:         true,
				StandardAnswer: false,
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
			var resQuestion QuestionJudgement
			err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wQuestion.Code)
			assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
			assert.Equal(b, question.Id, resQuestion.Id)
			assert.Equal(b, question.Title, resQuestion.Title)
			assert.Equal(b, question.Answer, resQuestion.Answer)
			assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
		}
	})
}

func TestNova_HandleCreateQuestionEssay(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateQuestionEssay (essay)
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
	url = server.URL + "/nova/v1/question/essay"
	question := QuestionEssay{
		Id:             reQuestionId,
		Title:          "What's the sweetest fruit?",
		Answer:         "apple",
		StandardAnswer: "apple",
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
	var resQuestion QuestionEssay
	err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wQuestion.Code)
	assert.Equal(t, "application/json", wQuestion.Header().Get("Content-Type"))
	assert.Equal(t, question.Id, resQuestion.Id)
	assert.Equal(t, question.Title, resQuestion.Title)
	assert.Equal(t, question.Answer, resQuestion.Answer)
	assert.Equal(t, question.StandardAnswer, resQuestion.StandardAnswer)
}

func BenchmarkNova_HandleCreateQuestionEssay(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestionEssay (essay)
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
		url = server.URL + "/nova/v1/question/essay"
		question := QuestionEssay{
			Id:             reQuestionId,
			Title:          "What's the sweetest fruit?",
			Answer:         "apple",
			StandardAnswer: "apple",
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
		var resQuestion QuestionEssay
		err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wQuestion.Code)
		assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
		assert.Equal(b, question.Id, resQuestion.Id)
		assert.Equal(b, question.Title, resQuestion.Title)
		assert.Equal(b, question.Answer, resQuestion.Answer)
		assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
	}
}

func BenchmarkNova_HandleCreateQuestionEssayParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestionEssay (essay) (Parallel)
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
			url = server.URL + "/nova/v1/question/essay"
			question := QuestionEssay{
				Id:             reQuestionId,
				Title:          "What's the sweetest fruit?",
				Answer:         "apple",
				StandardAnswer: "apple",
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
			var resQuestion QuestionEssay
			err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wQuestion.Code)
			assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
			assert.Equal(b, question.Id, resQuestion.Id)
			assert.Equal(b, question.Title, resQuestion.Title)
			assert.Equal(b, question.Answer, resQuestion.Answer)
			assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
		}
	})
}

func TestNova_HandleDeleteQuestionSingleChoice(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleCreateQuestion (single-choice)
	// Test Purpose: Test HandleCreateQuestion create question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
	url = server.URL + "/nova/v1/question/single-choice"
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
	/* delete question */
	// request content
	url = server.URL + "/nova/v1/question/single-choice"
	// request delete question
	wDeleteQuestion := httptest.NewRecorder()
	reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
	// validate response
	assert.Equal(t, http.StatusNoContent, wDeleteQuestion.Code)
}

func BenchmarkNova_HandleDeleteQuestionSingleChoice(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (single-choice)
	// Test Purpose: Benchmark HandleCreateQuestion create question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
		url = server.URL + "/nova/v1/question/single-choice"
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
		/* delete question */
		// request content
		url = server.URL + "/nova/v1/question/single-choice"
		// request delete question
		wDeleteQuestion := httptest.NewRecorder()
		reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
		// validate response
		assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
	}
}

func BenchmarkNova_HandleDeleteQuestionSingleChoiceParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleCreateQuestion (single-choice) (Parallel)
	// Test Purpose: Benchmark HandleCreateQuestion create question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
			url = server.URL + "/nova/v1/question/single-choice"
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
			/* delete question */
			// request content
			url = server.URL + "/nova/v1/question/single-choice"
			// request delete question
			wDeleteQuestion := httptest.NewRecorder()
			reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
			// validate response
			assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
		}
	})
}

func TestNova_HandleDeleteQuestionMultipleChoice(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleDeleteQuestionMultipleChoice (multiple-choice)
	// Test Purpose: Test HandleDeleteQuestionMultipleChoice delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
	url = server.URL + "/nova/v1/question/multiple-choice"
	question := QuestionMultipleChoice{
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
		StandardAnswers: []QuestionAnswer{
			{
				"B",
				"watermelon",
			},
			{
				"D",
				"peach",
			},
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
	var resQuestion QuestionMultipleChoice
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
	assert.Equal(t, question.StandardAnswers, resQuestion.StandardAnswers)
	/* delete question */
	// request content
	url = server.URL + "/nova/v1/question/multiple-choice"
	// request delete question
	wDeleteQuestion := httptest.NewRecorder()
	reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
	// validate response
	assert.Equal(t, http.StatusNoContent, wDeleteQuestion.Code)
}

func BenchmarkNova_HandleDeleteQuestionMultipleChoice(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteQuestionMultipleChoice (multiple-choice)
	// Test Purpose: Benchmark HandleDeleteQuestionMultipleChoice delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
		url = server.URL + "/nova/v1/question/multiple-choice"
		question := QuestionMultipleChoice{
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
			StandardAnswers: []QuestionAnswer{
				{
					"B",
					"watermelon",
				},
				{
					"D",
					"peach",
				},
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
		var resQuestion QuestionMultipleChoice
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
		assert.Equal(b, question.StandardAnswers, resQuestion.StandardAnswers)
		/* delete question */
		// request content
		url = server.URL + "/nova/v1/question/multiple-choice"
		// request delete question
		wDeleteQuestion := httptest.NewRecorder()
		reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
		// validate response
		assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
	}
}

func BenchmarkNova_HandleDeleteQuestionMultipleChoiceParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteQuestionMultipleChoice (multiple-choice)
	// Test Purpose: Benchmark HandleDeleteQuestionMultipleChoice delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
			url = server.URL + "/nova/v1/question/multiple-choice"
			question := QuestionMultipleChoice{
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
				StandardAnswers: []QuestionAnswer{
					{
						"B",
						"watermelon",
					},
					{
						"D",
						"peach",
					},
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
			var resQuestion QuestionMultipleChoice
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
			assert.Equal(b, question.StandardAnswers, resQuestion.StandardAnswers)
			/* delete question */
			// request content
			url = server.URL + "/nova/v1/question/multiple-choice"
			// request delete question
			wDeleteQuestion := httptest.NewRecorder()
			reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
			// validate response
			assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
		}
	})
}

func TestNova_HandleDeleteQuestionJudgement(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleDeleteQuestion (judgement)
	// Test Purpose: Test HandleDeleteQuestion delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
	url = server.URL + "/nova/v1/question/judgement"
	question := QuestionJudgement{
		Id:             reQuestionId,
		Title:          "What's the sweetest fruit?",
		Answer:         true,
		StandardAnswer: false,
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
	var resQuestion QuestionJudgement
	err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wQuestion.Code)
	assert.Equal(t, "application/json", wQuestion.Header().Get("Content-Type"))
	assert.Equal(t, question.Id, resQuestion.Id)
	assert.Equal(t, question.Title, resQuestion.Title)
	assert.Equal(t, question.Answer, resQuestion.Answer)
	assert.Equal(t, question.StandardAnswer, resQuestion.StandardAnswer)
	/* delete question */
	// request content
	url = server.URL + "/nova/v1/question/judgement"
	// request delete question
	wDeleteQuestion := httptest.NewRecorder()
	reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
	// validate response
	assert.Equal(t, http.StatusNoContent, wDeleteQuestion.Code)
}

func BenchmarkNova_HandleDeleteQuestionJudgement(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteQuestionEssay (judgement)
	// Test Purpose: Benchmark HandleDeleteQuestion delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
		url = server.URL + "/nova/v1/question/judgement"
		question := QuestionJudgement{
			Id:             reQuestionId,
			Title:          "What's the sweetest fruit?",
			Answer:         true,
			StandardAnswer: false,
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
		var resQuestion QuestionJudgement
		err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wQuestion.Code)
		assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
		assert.Equal(b, question.Id, resQuestion.Id)
		assert.Equal(b, question.Title, resQuestion.Title)
		assert.Equal(b, question.Answer, resQuestion.Answer)
		assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
		/* delete question */
		// request content
		url = server.URL + "/nova/v1/question/judgement"
		// request delete question
		wDeleteQuestion := httptest.NewRecorder()
		reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
		// validate response
		assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
	}
}

func BenchmarkNova_HandleDeleteQuestionJudgementParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteQuestionEssay (judgement) (parallel)
	// Test Purpose: Benchmark HandleDeleteQuestion delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
			url = server.URL + "/nova/v1/question/judgement"
			question := QuestionJudgement{
				Id:             reQuestionId,
				Title:          "What's the sweetest fruit?",
				Answer:         true,
				StandardAnswer: false,
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
			var resQuestion QuestionJudgement
			err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wQuestion.Code)
			assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
			assert.Equal(b, question.Id, resQuestion.Id)
			assert.Equal(b, question.Title, resQuestion.Title)
			assert.Equal(b, question.Answer, resQuestion.Answer)
			assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
			/* delete question */
			// request content
			url = server.URL + "/nova/v1/question/judgement"
			// request delete question
			wDeleteQuestion := httptest.NewRecorder()
			reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
			// validate response
			assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
		}
	})
}

func TestNova_HandleDeleteQuestionEssay(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleDeleteQuestion (essay)
	// Test Purpose: Test HandleDeleteQuestion delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
	url = server.URL + "/nova/v1/question/essay"
	question := QuestionEssay{
		Id:             reQuestionId,
		Title:          "What's the sweetest fruit?",
		Answer:         "apple",
		StandardAnswer: "apple",
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
	var resQuestion QuestionEssay
	err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wQuestion.Code)
	assert.Equal(t, "application/json", wQuestion.Header().Get("Content-Type"))
	assert.Equal(t, question.Id, resQuestion.Id)
	assert.Equal(t, question.Title, resQuestion.Title)
	assert.Equal(t, question.Answer, resQuestion.Answer)
	assert.Equal(t, question.StandardAnswer, resQuestion.StandardAnswer)
	/* delete question */
	// request content
	url = server.URL + "/nova/v1/question/essay"
	// request delete question
	wDeleteQuestion := httptest.NewRecorder()
	reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
	// validate response
	assert.Equal(t, http.StatusNoContent, wDeleteQuestion.Code)
}

func BenchmarkNova_HandleDeleteQuestionEssay(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteQuestionEssay (essay)
	// Test Purpose: Benchmark HandleDeleteQuestion delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
		url = server.URL + "/nova/v1/question/essay"
		question := QuestionEssay{
			Id:             reQuestionId,
			Title:          "What's the sweetest fruit?",
			Answer:         "apple",
			StandardAnswer: "apple",
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
		var resQuestion QuestionEssay
		err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wQuestion.Code)
		assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
		assert.Equal(b, question.Id, resQuestion.Id)
		assert.Equal(b, question.Title, resQuestion.Title)
		assert.Equal(b, question.Answer, resQuestion.Answer)
		assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
		/* delete question */
		// request content
		url = server.URL + "/nova/v1/question/essay"
		// request delete question
		wDeleteQuestion := httptest.NewRecorder()
		reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
		// validate response
		assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
	}
}

func BenchmarkNova_HandleDeleteQuestionEssayParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleDeleteQuestionEssay (essay) (parallel)
	// Test Purpose: Benchmark HandleDeleteQuestion delete question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send DeleteQuestion request with questionId by using DELETE method
	// 6. receive DeleteQuestion request by using 204 No Content Code
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
			url = server.URL + "/nova/v1/question/essay"
			question := QuestionEssay{
				Id:             reQuestionId,
				Title:          "What's the sweetest fruit?",
				Answer:         "apple",
				StandardAnswer: "apple",
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
			var resQuestion QuestionEssay
			err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wQuestion.Code)
			assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
			assert.Equal(b, question.Id, resQuestion.Id)
			assert.Equal(b, question.Title, resQuestion.Title)
			assert.Equal(b, question.Answer, resQuestion.Answer)
			assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
			/* delete question */
			// request content
			url = server.URL + "/nova/v1/question/essay"
			// request delete question
			wDeleteQuestion := httptest.NewRecorder()
			reqDeleteQuestion, err := http.NewRequest(http.MethodDelete, url+"/"+reQuestionId, nil)
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			router.ServeHTTP(wDeleteQuestion, reqDeleteQuestion)
			// validate response
			assert.Equal(b, http.StatusNoContent, wDeleteQuestion.Code)
		}
	})
}

func TestNova_HandleModifyQuestionSingleChoice(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleModifyQuestionSingleChoice (single-choice)
	// Test Purpose: Test HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
	url = server.URL + "/nova/v1/question/single-choice"
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
	/* modify question */
	// request content
	url = server.URL + "/nova/v1/question/single-choice"
	questionNew := QuestionSingleChoice{
		Id:    reQuestionId,
		Title: "Which city is the most one you favorite?",
		Answers: []QuestionAnswer{
			QuestionAnswer{
				"A",
				"New York",
			},
			QuestionAnswer{
				"B",
				"London",
			},
			QuestionAnswer{
				"C",
				"Paris",
			},
			QuestionAnswer{
				"D",
				"Beijing",
			},
		},
		StandardAnswer: QuestionAnswer{
			"D",
			"Beijing",
		},
	}
	bodyNew, err := json.Marshal(questionNew)
	if err != nil {
		t.Errorf("error marshal question: %v", err)
	}
	// request modify question
	wModifyQuestion := httptest.NewRecorder()
	reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqModifyQuestion.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
	// return response
	var resModifyQuestion QuestionSingleChoice
	err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
	// validate response
	assert.Equal(t, http.StatusOK, wModifyQuestion.Code)
	assert.Equal(t, "application/json", wModifyQuestion.Header().Get("Content-Type"))
	assert.Equal(t, questionNew.Id, resModifyQuestion.Id)
	assert.Equal(t, questionNew.Title, resModifyQuestion.Title)
	assert.Equal(t, questionNew.Answers, resModifyQuestion.Answers)
	assert.Equal(t, questionNew.StandardAnswer, resModifyQuestion.StandardAnswer)
}

func BenchmarkNova_HandleModifyQuestionSingleChoice(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleModifyQuestionSingleChoice (single-choice)
	// Test Purpose: Benchmark HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
		url = server.URL + "/nova/v1/question/single-choice"
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
		/* modify question */
		// request content
		url = server.URL + "/nova/v1/question/single-choice"
		questionNew := QuestionSingleChoice{
			Id:    reQuestionId,
			Title: "Which city is the most one you favorite?",
			Answers: []QuestionAnswer{
				QuestionAnswer{
					"A",
					"New York",
				},
				QuestionAnswer{
					"B",
					"London",
				},
				QuestionAnswer{
					"C",
					"Paris",
				},
				QuestionAnswer{
					"D",
					"Beijing",
				},
			},
			StandardAnswer: QuestionAnswer{
				"D",
				"Beijing",
			},
		}
		bodyNew, err := json.Marshal(questionNew)
		if err != nil {
			b.Errorf("error marshal question: %v", err)
		}
		// request modify question
		wModifyQuestion := httptest.NewRecorder()
		reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqModifyQuestion.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
		// return response
		var resModifyQuestion QuestionSingleChoice
		err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
		// validate response
		assert.Equal(b, http.StatusOK, wModifyQuestion.Code)
		assert.Equal(b, "application/json", wModifyQuestion.Header().Get("Content-Type"))
		assert.Equal(b, questionNew.Id, resModifyQuestion.Id)
		assert.Equal(b, questionNew.Title, resModifyQuestion.Title)
		assert.Equal(b, questionNew.Answers, resModifyQuestion.Answers)
		assert.Equal(b, questionNew.StandardAnswer, resModifyQuestion.StandardAnswer)
	}
}

func BenchmarkNova_HandleModifyQuestionSingleChoiceParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleModifyQuestionSingleChoice (single-choice) (parallel)
	// Test Purpose: Benchmark HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
			url = server.URL + "/nova/v1/question/single-choice"
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
			/* modify question */
			// request content
			url = server.URL + "/nova/v1/question/single-choice"
			questionNew := QuestionSingleChoice{
				Id:    reQuestionId,
				Title: "Which city is the most one you favorite?",
				Answers: []QuestionAnswer{
					QuestionAnswer{
						"A",
						"New York",
					},
					QuestionAnswer{
						"B",
						"London",
					},
					QuestionAnswer{
						"C",
						"Paris",
					},
					QuestionAnswer{
						"D",
						"Beijing",
					},
				},
				StandardAnswer: QuestionAnswer{
					"D",
					"Beijing",
				},
			}
			bodyNew, err := json.Marshal(questionNew)
			if err != nil {
				b.Errorf("error marshal question: %v", err)
			}
			// request modify question
			wModifyQuestion := httptest.NewRecorder()
			reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqModifyQuestion.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
			// return response
			var resModifyQuestion QuestionSingleChoice
			err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
			// validate response
			assert.Equal(b, http.StatusOK, wModifyQuestion.Code)
			assert.Equal(b, "application/json", wModifyQuestion.Header().Get("Content-Type"))
			assert.Equal(b, questionNew.Id, resModifyQuestion.Id)
			assert.Equal(b, questionNew.Title, resModifyQuestion.Title)
			assert.Equal(b, questionNew.Answers, resModifyQuestion.Answers)
			assert.Equal(b, questionNew.StandardAnswer, resModifyQuestion.StandardAnswer)
		}
	})
}

func TestNova_HandleModifyQuestionMultipleChoice(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleModifyQuestionMultipleChoice (multiple-choice)
	// Test Purpose: Test HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
	url = server.URL + "/nova/v1/question/multiple-choice"
	question := QuestionMultipleChoice{
		Id:    reQuestionId,
		Title: "What's the sweetest fruit?",
		Answers: []QuestionAnswer{
			{
				"A",
				"apple",
			},
			{
				"B",
				"watermelon",
			},
			{
				"C",
				"orange",
			},
			{
				"D",
				"peach",
			},
		},
		StandardAnswers: []QuestionAnswer{
			{
				"A",
				"apple",
			},
			{
				"B",
				"watermelon",
			},
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
	var resQuestion QuestionMultipleChoice
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
	assert.Equal(t, question.StandardAnswers, resQuestion.StandardAnswers)
	/* modify question */
	// request content
	url = server.URL + "/nova/v1/question/multiple-choice"
	questionNew := QuestionMultipleChoice{
		Id:    reQuestionId,
		Title: "Which city is the most one you favorite?",
		Answers: []QuestionAnswer{
			{
				"A",
				"New York",
			},
			{
				"B",
				"London",
			},
			{
				"C",
				"Paris",
			},
			{
				"D",
				"Beijing",
			},
		},
		StandardAnswers: []QuestionAnswer{
			{
				"C",
				"Paris",
			},
			{
				"D",
				"Beijing",
			},
		},
	}
	bodyNew, err := json.Marshal(questionNew)
	if err != nil {
		t.Errorf("error marshal question: %v", err)
	}
	// request modify question
	wModifyQuestion := httptest.NewRecorder()
	reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqModifyQuestion.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
	// return response
	var resModifyQuestion QuestionMultipleChoice
	err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
	// validate response
	assert.Equal(t, http.StatusOK, wModifyQuestion.Code)
	assert.Equal(t, "application/json", wModifyQuestion.Header().Get("Content-Type"))
	assert.Equal(t, questionNew.Id, resModifyQuestion.Id)
	assert.Equal(t, questionNew.Title, resModifyQuestion.Title)
	assert.Equal(t, questionNew.Answers, resModifyQuestion.Answers)
	assert.Equal(t, questionNew.StandardAnswers, resModifyQuestion.StandardAnswers)
}

func BenchmarkNova_HandleModifyQuestionMultipleChoice(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleModifyQuestionMultipleChoice (multiple-choice)
	// Test Purpose: Benchmark HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
		url = server.URL + "/nova/v1/question/multiple-choice"
		question := QuestionMultipleChoice{
			Id:    reQuestionId,
			Title: "What's the sweetest fruit?",
			Answers: []QuestionAnswer{
				{
					"A",
					"apple",
				},
				{
					"B",
					"watermelon",
				},
				{
					"C",
					"orange",
				},
				{
					"D",
					"peach",
				},
			},
			StandardAnswers: []QuestionAnswer{
				{
					"A",
					"apple",
				},
				{
					"B",
					"watermelon",
				},
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
		var resQuestion QuestionMultipleChoice
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
		assert.Equal(b, question.StandardAnswers, resQuestion.StandardAnswers)
		/* modify question */
		// request content
		url = server.URL + "/nova/v1/question/multiple-choice"
		questionNew := QuestionMultipleChoice{
			Id:    reQuestionId,
			Title: "Which city is the most one you favorite?",
			Answers: []QuestionAnswer{
				{
					"A",
					"New York",
				},
				{
					"B",
					"London",
				},
				{
					"C",
					"Paris",
				},
				{
					"D",
					"Beijing",
				},
			},
			StandardAnswers: []QuestionAnswer{
				{
					"C",
					"Paris",
				},
				{
					"D",
					"Beijing",
				},
			},
		}
		bodyNew, err := json.Marshal(questionNew)
		if err != nil {
			b.Errorf("error marshal question: %v", err)
		}
		// request modify question
		wModifyQuestion := httptest.NewRecorder()
		reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqModifyQuestion.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
		// return response
		var resModifyQuestion QuestionMultipleChoice
		err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
		// validate response
		assert.Equal(b, http.StatusOK, wModifyQuestion.Code)
		assert.Equal(b, "application/json", wModifyQuestion.Header().Get("Content-Type"))
		assert.Equal(b, questionNew.Id, resModifyQuestion.Id)
		assert.Equal(b, questionNew.Title, resModifyQuestion.Title)
		assert.Equal(b, questionNew.Answers, resModifyQuestion.Answers)
		assert.Equal(b, questionNew.StandardAnswers, resModifyQuestion.StandardAnswers)
	}
}

func BenchmarkNova_HandleModifyQuestionMultipleChoiceParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleModifyQuestionMultipleChoice (multiple-choice) (parallel)
	// Test Purpose: Benchmark HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
			url = server.URL + "/nova/v1/question/multiple-choice"
			question := QuestionMultipleChoice{
				Id:    reQuestionId,
				Title: "What's the sweetest fruit?",
				Answers: []QuestionAnswer{
					{
						"A",
						"apple",
					},
					{
						"B",
						"watermelon",
					},
					{
						"C",
						"orange",
					},
					{
						"D",
						"peach",
					},
				},
				StandardAnswers: []QuestionAnswer{
					{
						"A",
						"apple",
					},
					{
						"B",
						"watermelon",
					},
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
			var resQuestion QuestionMultipleChoice
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
			assert.Equal(b, question.StandardAnswers, resQuestion.StandardAnswers)
			/* modify question */
			// request content
			url = server.URL + "/nova/v1/question/multiple-choice"
			questionNew := QuestionMultipleChoice{
				Id:    reQuestionId,
				Title: "Which city is the most one you favorite?",
				Answers: []QuestionAnswer{
					{
						"A",
						"New York",
					},
					{
						"B",
						"London",
					},
					{
						"C",
						"Paris",
					},
					{
						"D",
						"Beijing",
					},
				},
				StandardAnswers: []QuestionAnswer{
					{
						"C",
						"Paris",
					},
					{
						"D",
						"Beijing",
					},
				},
			}
			bodyNew, err := json.Marshal(questionNew)
			if err != nil {
				b.Errorf("error marshal question: %v", err)
			}
			// request modify question
			wModifyQuestion := httptest.NewRecorder()
			reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqModifyQuestion.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
			// return response
			var resModifyQuestion QuestionMultipleChoice
			err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
			// validate response
			assert.Equal(b, http.StatusOK, wModifyQuestion.Code)
			assert.Equal(b, "application/json", wModifyQuestion.Header().Get("Content-Type"))
			assert.Equal(b, questionNew.Id, resModifyQuestion.Id)
			assert.Equal(b, questionNew.Title, resModifyQuestion.Title)
			assert.Equal(b, questionNew.Answers, resModifyQuestion.Answers)
			assert.Equal(b, questionNew.StandardAnswers, resModifyQuestion.StandardAnswers)
		}
	})
}

func TestNova_HandleModifyQuestionJudgement(t *testing.T) {
	/*---------------------------------------------------------------------------------------
	// Test Case: TestNova_HandleModifyQuestionJudgement (judgement)
	// Test Purpose: Test HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
	url = server.URL + "/nova/v1/question/judgement"
	question := QuestionJudgement{
		Id:             reQuestionId,
		Title:          "What's the sweetest fruit?",
		Answer:         true,
		StandardAnswer: false,
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
	var resQuestion QuestionJudgement
	err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
	if err != nil {
		t.Errorf("error unmarshal response: %v", err)
	}
	// validate response
	assert.Equal(t, http.StatusCreated, wQuestion.Code)
	assert.Equal(t, "application/json", wQuestion.Header().Get("Content-Type"))
	assert.Equal(t, question.Id, resQuestion.Id)
	assert.Equal(t, question.Title, resQuestion.Title)
	assert.Equal(t, question.Answer, resQuestion.Answer)
	assert.Equal(t, question.StandardAnswer, resQuestion.StandardAnswer)
	/* modify question */
	// request content
	url = server.URL + "/nova/v1/question/judgement"
	questionNew := QuestionJudgement{
		Id:             reQuestionId,
		Title:          "Is New York in America?",
		Answer:         false,
		StandardAnswer: true,
	}
	bodyNew, err := json.Marshal(questionNew)
	if err != nil {
		t.Errorf("error marshal question: %v", err)
	}
	// request modify question
	wModifyQuestion := httptest.NewRecorder()
	reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}
	reqModifyQuestion.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
	// return response
	var resModifyQuestion QuestionJudgement
	err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
	// validate response
	assert.Equal(t, http.StatusOK, wModifyQuestion.Code)
	assert.Equal(t, "application/json", wModifyQuestion.Header().Get("Content-Type"))
	assert.Equal(t, questionNew.Id, resModifyQuestion.Id)
	assert.Equal(t, questionNew.Title, resModifyQuestion.Title)
	assert.Equal(t, questionNew.Answer, resModifyQuestion.Answer)
	assert.Equal(t, questionNew.StandardAnswer, resModifyQuestion.StandardAnswer)
}

func BenchmarkNova_HandleModifyQuestionJudgement(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleModifyQuestionJudgement (judgement)
	// Test Purpose: Benchmark HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
		url = server.URL + "/nova/v1/question/judgement"
		question := QuestionJudgement{
			Id:             reQuestionId,
			Title:          "What's the sweetest fruit?",
			Answer:         true,
			StandardAnswer: false,
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
		var resQuestion QuestionJudgement
		err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
		if err != nil {
			b.Errorf("error unmarshal response: %v", err)
		}
		// validate response
		assert.Equal(b, http.StatusCreated, wQuestion.Code)
		assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
		assert.Equal(b, question.Id, resQuestion.Id)
		assert.Equal(b, question.Title, resQuestion.Title)
		assert.Equal(b, question.Answer, resQuestion.Answer)
		assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
		/* modify question */
		// request content
		url = server.URL + "/nova/v1/question/judgement"
		questionNew := QuestionJudgement{
			Id:             reQuestionId,
			Title:          "Is New York in America?",
			Answer:         false,
			StandardAnswer: true,
		}
		bodyNew, err := json.Marshal(questionNew)
		if err != nil {
			b.Errorf("error marshal question: %v", err)
		}
		// request modify question
		wModifyQuestion := httptest.NewRecorder()
		reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
		if err != nil {
			b.Errorf("error creating request: %v", err)
		}
		reqModifyQuestion.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
		// return response
		var resModifyQuestion QuestionJudgement
		err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
		// validate response
		assert.Equal(b, http.StatusOK, wModifyQuestion.Code)
		assert.Equal(b, "application/json", wModifyQuestion.Header().Get("Content-Type"))
		assert.Equal(b, questionNew.Id, resModifyQuestion.Id)
		assert.Equal(b, questionNew.Title, resModifyQuestion.Title)
		assert.Equal(b, questionNew.Answer, resModifyQuestion.Answer)
		assert.Equal(b, questionNew.StandardAnswer, resModifyQuestion.StandardAnswer)
	}
}

func BenchmarkNova_HandleModifyQuestionJudgementParallel(b *testing.B) {
	/*---------------------------------------------------------------------------------------
	// Test Case: BenchmarkNova_HandleModifyQuestionJudgement (judgement) (parallel)
	// Test Purpose: Benchmark HandleModifyQuestion modify question
	// Test Steps:
	// 1. send CreateQuestionId request by using POST method
	// 2. receive CreateQuestionId response with created questionId by using 201 Created Code
	// 3. send CreateQuestion request by using POST method
	// 4. receive CreateQuestion response with created question by using 201 Created Code
	// 5. send ModifyQuestion request with questionId by using PATCH method
	// 6. receive ModifyQuestion request by using 200 No Content Code
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
			url = server.URL + "/nova/v1/question/judgement"
			question := QuestionJudgement{
				Id:             reQuestionId,
				Title:          "What's the sweetest fruit?",
				Answer:         true,
				StandardAnswer: false,
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
			var resQuestion QuestionJudgement
			err = json.Unmarshal(wQuestion.Body.Bytes(), &resQuestion)
			if err != nil {
				b.Errorf("error unmarshal response: %v", err)
			}
			// validate response
			assert.Equal(b, http.StatusCreated, wQuestion.Code)
			assert.Equal(b, "application/json", wQuestion.Header().Get("Content-Type"))
			assert.Equal(b, question.Id, resQuestion.Id)
			assert.Equal(b, question.Title, resQuestion.Title)
			assert.Equal(b, question.Answer, resQuestion.Answer)
			assert.Equal(b, question.StandardAnswer, resQuestion.StandardAnswer)
			/* modify question */
			// request content
			url = server.URL + "/nova/v1/question/judgement"
			questionNew := QuestionJudgement{
				Id:             reQuestionId,
				Title:          "Is New York in America?",
				Answer:         false,
				StandardAnswer: true,
			}
			bodyNew, err := json.Marshal(questionNew)
			if err != nil {
				b.Errorf("error marshal question: %v", err)
			}
			// request modify question
			wModifyQuestion := httptest.NewRecorder()
			reqModifyQuestion, err := http.NewRequest(http.MethodPatch, url+"/"+reQuestionId, bytes.NewReader(bodyNew))
			if err != nil {
				b.Errorf("error creating request: %v", err)
			}
			reqModifyQuestion.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(wModifyQuestion, reqModifyQuestion)
			// return response
			var resModifyQuestion QuestionJudgement
			err = json.Unmarshal(wModifyQuestion.Body.Bytes(), &resModifyQuestion)
			// validate response
			assert.Equal(b, http.StatusOK, wModifyQuestion.Code)
			assert.Equal(b, "application/json", wModifyQuestion.Header().Get("Content-Type"))
			assert.Equal(b, questionNew.Id, resModifyQuestion.Id)
			assert.Equal(b, questionNew.Title, resModifyQuestion.Title)
			assert.Equal(b, questionNew.Answer, resModifyQuestion.Answer)
			assert.Equal(b, questionNew.StandardAnswer, resModifyQuestion.StandardAnswer)
		}
	})
}

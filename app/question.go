package app

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nova/logger"
	"reflect"
	"strings"
)

func (nova *Nova) HandleCreateQuestionId(c *gin.Context) {
	// create questionId
	var questionId string
	logger.Infof("handle request create questionId")
	// generate questionId
	questionId = uuid.New().String()
	logger.Debugf("generate questionId: %v", questionId)
	// return response
	nova.response201Created(c, questionId)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, questionId)
	return
}

func (nova *Nova) HandleCreateQuestion(c *gin.Context) {
	// create question
	logger.Infof("handle request create question")
	// get raw data from request body
	logger.Debugf("get raw data from request body")
	rawData, err := c.GetRawData()
	if err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error get raw data from request body: %v", err)
		return
	}
	logger.Debugf("successfully get raw data from request body")
	// reflect request type
	logger.Debugf("reflect request type")
	reqType, err := nova.reflectRequestType(rawData)
	if err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error reflect request type: %v", err)
		return
	}
	logger.Debugf("successfully reflect request type")
	// binding structure according request type
	switch reqType {
	case "single_choice":
		nova.handleCreateQuestionSingleChoice(c, rawData)
	case "multiple_choice":
		nova.handleCreateQuestionMultipleChoice(c, rawData)
	case "judgement":
		nova.handleCreateQuestionJudgement(c, rawData)
	case "essay":
		nova.handleCreateQuestionEssay(c, rawData)
	default:
		nova.response400BadRequest(c, errors.New("invalid request type"))
	}
	return
}

func (nova *Nova) handleCreateQuestionSingleChoice(c *gin.Context, rawData []byte) {
	// create question
	var request QuestionSingleChoice
	logger.Infof("handle request create single-choice question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check single-choice question is validate")
	b, err := nova.isQuestionSingleChoiceValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check single-choice question is validate: %v", err)
		return
	}
	logger.Debugf("successfully check single-choice question is validate")
	// update data cache by querying single-choice questions in database
	logger.Debugf("update data cache by querying single-choice question in database")
	err = nova.querySingleChoiceQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying single-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying single-choice question in database")
	// check user existence
	logger.Debugf("check single-choice question is existed")
	if nova.isSingleChoiceQuestionExisted(strings.ToLower(request.Id)) {
		nova.response409Conflict(c, errors.New("single-choice question already exists"))
		logger.Errorf("error check single-choice question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check single-choice question is existed")
	// store created single-choice question in data cache
	logger.Debugf("store single-choice question in data cache")
	response := QuestionSingleChoice{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answers:        request.Answers,
		StandardAnswer: request.StandardAnswer,
	}
	nova.createSingleChoiceQuestionInDataCache(response)
	logger.Debugf("successfully store single-choice question in data cache")
	// store created single-choice question in database
	logger.Debugf("store single-choice question in database")
	if err = nova.createSingleChoiceQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error single-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully store single-choice question in database")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) handleCreateQuestionMultipleChoice(c *gin.Context, rawData []byte) {
	// create question
	var request QuestionMultipleChoice
	logger.Infof("handle request create question multiple-choice")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check question multiple-choice is validate")
	b, err := nova.isQuestionMultipleChoiceValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check question multiple-choice is validate: %v", err)
		return
	}
	logger.Debugf("successfully check question multiple-choice is validate")
	// store created question in data cache
	logger.Debugf("store question in data cache")
	response := QuestionMultipleChoice{
		Id:              strings.ToLower(request.Id),
		Title:           request.Title,
		Answers:         request.Answers,
		StandardAnswers: request.StandardAnswers,
	}
	// nova.createUserInDataCache(response)
	logger.Debugf("successfully store question in data cache")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) handleCreateQuestionJudgement(c *gin.Context, rawData []byte) {
	// create question
	var request QuestionJudgement
	logger.Infof("handle request create question judgement")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check question judgement is validate")
	b, err := nova.isQuestionJudgementValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check question judgement is validate: %v", err)
		return
	}
	logger.Debugf("successfully check question judgement is validate")
	// store created question in data cache
	logger.Debugf("store question in data cache")
	response := QuestionJudgement{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answer:         request.Answer,
		StandardAnswer: request.StandardAnswer,
	}
	// nova.createUserInDataCache(response)
	logger.Debugf("successfully store question in data cache")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) handleCreateQuestionEssay(c *gin.Context, rawData []byte) {
	// create question
	var request QuestionEssay
	logger.Infof("handle request create question judgement")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check question essay is validate")
	b, err := nova.isQuestionEssayValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check question essay is validate: %v", err)
		return
	}
	logger.Debugf("successfully check question essay is validate")
	// store created question in data cache
	logger.Debugf("store question in data cache")
	response := QuestionEssay{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answer:         request.Answer,
		StandardAnswer: request.StandardAnswer,
	}
	// nova.createUserInDataCache(response)
	logger.Debugf("successfully store question in data cache")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) HandleDeleteQuestion(c *gin.Context) {

}

func (nova *Nova) HandleModifyQuestion(c *gin.Context) {

}

func (nova *Nova) HandleQueryQuestion(c *gin.Context) {

}

func (nova *Nova) HandleUpdateQuestion(c *gin.Context) {

}

func (nova *Nova) reflectRequestType(rawData []byte) (string, error) {
	var temp map[string]interface{}
	// unmarshal raw data to temp map interface
	if err := json.Unmarshal(rawData, &temp); err != nil {
		return "", err
	}
	// check mandatory segment existence
	if nova.isRequiredFields(temp, QuestionSingleChoice{}) {
		return "single_choice", nil
	}
	if nova.isRequiredFields(temp, QuestionMultipleChoice{}) {
		return "multiple_choice", nil
	}
	if nova.isRequiredFields(temp, QuestionJudgement{}) {
		return "judgement", nil
	}
	if nova.isRequiredFields(temp, QuestionEssay{}) {
		return "essay", nil
	}
	return "", errors.New("invalid request")
}

func (nova *Nova) isRequiredFields(data map[string]interface{}, model interface{}) bool {
	// reflect type of model
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		binding := field.Tag.Get("binding")
		if tag != "" && nova.containsBinding(binding, "required") {
			if _, exists := data[tag]; !exists {
				return false
			}
		}
	}
	return true
}

func (nova *Nova) containsBinding(binding, rule string) bool {
	if binding == "" {
		return false
	}
	// split binding rules
	for _, r := range nova.splitBindingRules(binding) {
		if r == rule {
			return true
		}
	}
	return false
}

func (nova *Nova) splitBindingRules(binding string) []string {
	var rules []string
	start := 0
	// check binding char
	for i, char := range binding {
		if char == ',' {
			rules = append(rules, binding[start:i])
			start = i + 1
		}
	}
	// check binding length
	if start < len(binding) {
		rules = append(rules, binding[start:])
	}
	return rules
}

func (nova *Nova) isSingleChoiceQuestionExisted(id string) bool {
	// enable single-choice question cache read lock
	nova.cache.questionsCache.singleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.RUnlock()
	// search Id in data cache
	for _, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == id {
			return true
		}
	}
	return false
}

func (nova *Nova) isQuestionSingleChoiceValidate(question QuestionSingleChoice) (bool, error) {
	// check question identity format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) isQuestionMultipleChoiceValidate(question QuestionMultipleChoice) (bool, error) {
	// check question identity format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) isQuestionJudgementValidate(question QuestionJudgement) (bool, error) {
	// check question identity format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) isQuestionEssayValidate(question QuestionEssay) (bool, error) {
	// check question identity format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) createSingleChoiceQuestionInDataCache(question QuestionSingleChoice) {
	// enable single-choice question cache write lock
	nova.cache.questionsCache.singleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.Unlock()
	// append single-choice question in data cache
	nova.cache.questionsCache.singleChoiceCache.singleChoiceSet = append(nova.cache.questionsCache.singleChoiceCache.singleChoiceSet, question)
	return
}

func (nova *Nova) querySingleChoiceQuestionsInDatabase() error {
	// enable single-choice question cache write lock
	nova.cache.questionsCache.singleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.Unlock()
	// query single-choice question from database
	questions, err := nova.db.QueryQuestionsSingleChoice()
	if err != nil {
		return err
	}
	// update single-choice question in data cache
	for _, question := range questions {
		b := false
		// update if single-choice question existed
		for k, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
			if v.Id == question.Id {
				nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[k] = *question
				b = true
				break
			}
		}
		// create single-choice question if question not existed
		if !b {
			nova.cache.questionsCache.singleChoiceCache.singleChoiceSet = append(nova.cache.questionsCache.singleChoiceCache.singleChoiceSet, *question)
		}
	}
	return nil
}

func (nova *Nova) createSingleChoiceQuestionInDatabase(id string) error {
	// enable single-choice question cache read lock
	nova.cache.questionsCache.singleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.RUnlock()
	// search Id in data cache
	b := false
	question := QuestionSingleChoice{}
	for _, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == id {
			question = v
			b = true
			break
		}
	}
	if !b {
		return errors.New("single-choice question not found")
	}
	// create single-choice question in database
	if _, err := nova.db.CreateQuestionSingleChoice(&question); err != nil {
		return err
	}
	return nil
}

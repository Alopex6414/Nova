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
	logger.Infof("handle request create question single-choice")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check question single-choice is validate")
	b, err := nova.isQuestionSingleChoiceValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check question single-choice is validate: %v", err)
		return
	}
	logger.Debugf("successfully check question single-choice is validate")
	// store created question in data cache
	logger.Debugf("store question in data cache")
	response := QuestionSingleChoice{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answers:        request.Answers,
		StandardAnswer: request.StandardAnswer,
	}
	// nova.createUserInDataCache(response)
	logger.Debugf("successfully store question in data cache")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) handleCreateQuestionMultipleChoice(c *gin.Context, rawData []byte) {
	return
}

func (nova *Nova) handleCreateQuestionJudgement(c *gin.Context, rawData []byte) {
	return
}

func (nova *Nova) handleCreateQuestionEssay(c *gin.Context, rawData []byte) {
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

func (nova *Nova) isQuestionSingleChoiceValidate(question QuestionSingleChoice) (bool, error) {
	// check question Id format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

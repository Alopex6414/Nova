package app

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nova/logger"
	"reflect"
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
		nova.HandleCreateQuestionSingleChoice(c)
	case "multiple_choice":
		nova.HandleCreateQuestionMultipleChoice(c)
	case "judgement":
		nova.HandleCreateQuestionJudgement(c)
	case "essay":
		nova.HandleCreateQuestionEssay(c)
	default:
		nova.response400BadRequest(c, errors.New("invalid request type"))
	}
	return
}

func (nova *Nova) HandleCreateQuestionSingleChoice(c *gin.Context) {
	return
}

func (nova *Nova) HandleCreateQuestionMultipleChoice(c *gin.Context) {
	return
}

func (nova *Nova) HandleCreateQuestionJudgement(c *gin.Context) {
	return
}

func (nova *Nova) HandleCreateQuestionEssay(c *gin.Context) {
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

func (nova *Nova) reflectRequestType(raw []byte) (string, error) {
	var data map[string]interface{}
	// unmarshal raw data to data map interface
	if err := json.Unmarshal(raw, &data); err != nil {
		return "", err
	}
	// check mandatory segment existence
	if nova.isRequiredFields(data, QuestionSingleChoice{}) {
		return "single_choice", nil
	}
	if nova.isRequiredFields(data, QuestionMultipleChoice{}) {
		return "multiple_choice", nil
	}
	if nova.isRequiredFields(data, QuestionJudgement{}) {
		return "judgement", nil
	}
	if nova.isRequiredFields(data, QuestionEssay{}) {
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

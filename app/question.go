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
	// create single-choice question
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
	b, err := nova.isSingleChoiceQuestionValidate(request)
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
	// check single-choice question existence
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
	// create multiple-choice question
	var request QuestionMultipleChoice
	logger.Infof("handle request create multiple-choice question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check multiple-choice question is validate")
	b, err := nova.isMultipleChoiceQuestionValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check multiple-choice question is validate: %v", err)
		return
	}
	logger.Debugf("successfully check multiple-choice question is validate")
	// update data cache by querying multiple-choice questions in database
	logger.Debugf("update data cache by querying multiple-choice question in database")
	err = nova.queryMultipleChoiceQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying multiple-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying multiple-choice question in database")
	// check multiple-choice question existence
	logger.Debugf("check multiple-choice question is existed")
	if nova.isMultipleChoiceQuestionExisted(strings.ToLower(request.Id)) {
		nova.response409Conflict(c, errors.New("multiple-choice question already exists"))
		logger.Errorf("error check multiple-choice question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check multiple-choice question is existed")
	// store created multiple-choice question in data cache
	logger.Debugf("store multiple-choice question in data cache")
	response := QuestionMultipleChoice{
		Id:              strings.ToLower(request.Id),
		Title:           request.Title,
		Answers:         request.Answers,
		StandardAnswers: request.StandardAnswers,
	}
	nova.createMultipleChoiceQuestionInDataCache(response)
	logger.Debugf("successfully store multiple-choice question in data cache")
	// store created multiple-choice question in database
	logger.Debugf("store multiple-choice question in database")
	if err = nova.createMultipleChoiceQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error multiple-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully store multiple-choice question in database")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) handleCreateQuestionJudgement(c *gin.Context, rawData []byte) {
	// create judgement question
	var request QuestionJudgement
	logger.Infof("handle request create judgement question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check judgement question is validate")
	b, err := nova.isJudgementQuestionValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check judgement question is validate: %v", err)
		return
	}
	logger.Debugf("successfully check judgement question is validate")
	// update data cache by querying judgement questions in database
	logger.Debugf("update data cache by querying judgement question in database")
	err = nova.queryJudgementQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying judgement question in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying judgement question in database")
	// check judgement question existence
	logger.Debugf("check judgement question is existed")
	if nova.isJudgementQuestionExisted(strings.ToLower(request.Id)) {
		nova.response409Conflict(c, errors.New("judgement question already exists"))
		logger.Errorf("error check judgement question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check judgement question is existed")
	// store created judgement question in data cache
	logger.Debugf("store judgement question in data cache")
	response := QuestionJudgement{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answer:         request.Answer,
		StandardAnswer: request.StandardAnswer,
	}
	nova.createJudgementQuestionInDataCache(response)
	logger.Debugf("successfully store judgement question in data cache")
	// store created judgement question in database
	logger.Debugf("store judgement question in database")
	if err = nova.createJudgementQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error judgement question in database: %v", err)
		return
	}
	logger.Debugf("successfully store judgement question in database")
	// return response
	nova.response201Created(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusCreated, response)
	return
}

func (nova *Nova) handleCreateQuestionEssay(c *gin.Context, rawData []byte) {
	// create essay question
	var request QuestionEssay
	logger.Infof("handle request create essay question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := json.Unmarshal(rawData, &request); err != nil {
		nova.response400BadRequest(c, err)
		logger.Errorf("error bind request to json: %v", err)
		return
	}
	logger.Debugf("successfully bind request json format")
	// check request body correctness
	logger.Debugf("check essay question is validate")
	b, err := nova.isEssayQuestionValidate(request)
	if !b {
		nova.response400BadRequest(c, err)
		logger.Errorf("error check essay question is validate: %v", err)
		return
	}
	logger.Debugf("successfully check essay question is validate")
	// update data cache by querying essay questions in database
	logger.Debugf("update data cache by querying essay question in database")
	err = nova.queryEssayQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying essay question in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying essay question in database")
	// check essay question existence
	logger.Debugf("check essay question is existed")
	if nova.isEssayQuestionExisted(strings.ToLower(request.Id)) {
		nova.response409Conflict(c, errors.New("essay question already exists"))
		logger.Errorf("error check essay question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check essay question is existed")
	// store created essay question in data cache
	logger.Debugf("store judgement question in data cache")
	response := QuestionEssay{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answer:         request.Answer,
		StandardAnswer: request.StandardAnswer,
	}
	nova.createEssayQuestionInDataCache(response)
	logger.Debugf("successfully store essay question in data cache")
	// store created essay question in database
	logger.Debugf("store essay question in database")
	if err = nova.createEssayQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error essay question in database: %v", err)
		return
	}
	logger.Debugf("successfully store essay question in database")
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

func (nova *Nova) isMultipleChoiceQuestionExisted(id string) bool {
	// enable multiple-choice question cache read lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.RUnlock()
	// search Id in data cache
	for _, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == id {
			return true
		}
	}
	return false
}

func (nova *Nova) isJudgementQuestionExisted(id string) bool {
	// enable judgement question cache read lock
	nova.cache.questionsCache.judgementCache.mutex.RLock()
	defer nova.cache.questionsCache.judgementCache.mutex.RUnlock()
	// search Id in data cache
	for _, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == id {
			return true
		}
	}
	return false
}

func (nova *Nova) isEssayQuestionExisted(id string) bool {
	// enable essay question cache read lock
	nova.cache.questionsCache.essayCache.mutex.RLock()
	defer nova.cache.questionsCache.essayCache.mutex.RUnlock()
	// search Id in data cache
	for _, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == id {
			return true
		}
	}
	return false
}

func (nova *Nova) isSingleChoiceQuestionValidate(question QuestionSingleChoice) (bool, error) {
	// check question identity format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) isMultipleChoiceQuestionValidate(question QuestionMultipleChoice) (bool, error) {
	// check question identity format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) isJudgementQuestionValidate(question QuestionJudgement) (bool, error) {
	// check question identity format is UUID
	err := uuid.Validate(question.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (nova *Nova) isEssayQuestionValidate(question QuestionEssay) (bool, error) {
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

func (nova *Nova) createMultipleChoiceQuestionInDataCache(question QuestionMultipleChoice) {
	// enable multiple-choice question cache write lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.Unlock()
	// append multiple-choice question in data cache
	nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet = append(nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet, question)
	return
}

func (nova *Nova) queryMultipleChoiceQuestionsInDatabase() error {
	// enable multiple-choice question cache write lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.Unlock()
	// query multiple-choice question from database
	questions, err := nova.db.QueryQuestionsMultipleChoice()
	if err != nil {
		return err
	}
	// update multiple-choice question in data cache
	for _, question := range questions {
		b := false
		// update if multiple-choice question existed
		for k, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
			if v.Id == question.Id {
				nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[k] = *question
				b = true
				break
			}
		}
		// create multiple-choice question if question not existed
		if !b {
			nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet = append(nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet, *question)
		}
	}
	return nil
}

func (nova *Nova) createMultipleChoiceQuestionInDatabase(id string) error {
	// enable multiple-choice question cache read lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.RUnlock()
	// search Id in data cache
	b := false
	question := QuestionMultipleChoice{}
	for _, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == id {
			question = v
			b = true
			break
		}
	}
	if !b {
		return errors.New("multiple-choice question not found")
	}
	// create multiple-choice question in database
	if _, err := nova.db.CreateQuestionMultipleChoice(&question); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) createJudgementQuestionInDataCache(question QuestionJudgement) {
	// enable judgement question cache write lock
	nova.cache.questionsCache.judgementCache.mutex.Lock()
	defer nova.cache.questionsCache.judgementCache.mutex.Unlock()
	// append judgement question in data cache
	nova.cache.questionsCache.judgementCache.judgementSet = append(nova.cache.questionsCache.judgementCache.judgementSet, question)
	return
}

func (nova *Nova) queryJudgementQuestionsInDatabase() error {
	// enable judgement question cache write lock
	nova.cache.questionsCache.judgementCache.mutex.Lock()
	defer nova.cache.questionsCache.judgementCache.mutex.Unlock()
	// query judgement question from database
	questions, err := nova.db.QueryQuestionsJudgement()
	if err != nil {
		return err
	}
	// update judgement question in data cache
	for _, question := range questions {
		b := false
		// update if judgement question existed
		for k, v := range nova.cache.questionsCache.judgementCache.judgementSet {
			if v.Id == question.Id {
				nova.cache.questionsCache.judgementCache.judgementSet[k] = *question
				b = true
				break
			}
		}
		// create judgement question if question not existed
		if !b {
			nova.cache.questionsCache.judgementCache.judgementSet = append(nova.cache.questionsCache.judgementCache.judgementSet, *question)
		}
	}
	return nil
}

func (nova *Nova) createJudgementQuestionInDatabase(id string) error {
	// enable judgement question cache read lock
	nova.cache.questionsCache.judgementCache.mutex.RLock()
	defer nova.cache.questionsCache.judgementCache.mutex.RUnlock()
	// search Id in data cache
	b := false
	question := QuestionJudgement{}
	for _, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == id {
			question = v
			b = true
			break
		}
	}
	if !b {
		return errors.New("judgement question not found")
	}
	// create judgement question in database
	if _, err := nova.db.CreateQuestionJudgement(&question); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) createEssayQuestionInDataCache(question QuestionEssay) {
	// enable essay question cache write lock
	nova.cache.questionsCache.essayCache.mutex.Lock()
	defer nova.cache.questionsCache.essayCache.mutex.Unlock()
	// append essay question in data cache
	nova.cache.questionsCache.essayCache.essaySet = append(nova.cache.questionsCache.essayCache.essaySet, question)
	return
}

func (nova *Nova) queryEssayQuestionsInDatabase() error {
	// enable essay question cache write lock
	nova.cache.questionsCache.essayCache.mutex.Lock()
	defer nova.cache.questionsCache.essayCache.mutex.Unlock()
	// query essay question from database
	questions, err := nova.db.QueryQuestionsEssay()
	if err != nil {
		return err
	}
	// update essay question in data cache
	for _, question := range questions {
		b := false
		// update if essay question existed
		for k, v := range nova.cache.questionsCache.essayCache.essaySet {
			if v.Id == question.Id {
				nova.cache.questionsCache.essayCache.essaySet[k] = *question
				b = true
				break
			}
		}
		// create essay question if question not existed
		if !b {
			nova.cache.questionsCache.essayCache.essaySet = append(nova.cache.questionsCache.essayCache.essaySet, *question)
		}
	}
	return nil
}

func (nova *Nova) createEssayQuestionInDatabase(id string) error {
	// enable essay question cache read lock
	nova.cache.questionsCache.essayCache.mutex.RLock()
	defer nova.cache.questionsCache.essayCache.mutex.RUnlock()
	// search Id in data cache
	b := false
	question := QuestionEssay{}
	for _, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == id {
			question = v
			b = true
			break
		}
	}
	if !b {
		return errors.New("essay question not found")
	}
	// create essay question in database
	if _, err := nova.db.CreateQuestionEssay(&question); err != nil {
		return err
	}
	return nil
}

package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nova/logger"
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

func (nova *Nova) HandleCreateQuestionSingleChoice(c *gin.Context) {
	// create single-choice question
	var request QuestionSingleChoice
	logger.Infof("handle request create single-choice question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.ShouldBindJSON(&request); err != nil {
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

func (nova *Nova) HandleCreateQuestionMultipleChoice(c *gin.Context) {
	// create multiple-choice question
	var request QuestionMultipleChoice
	logger.Infof("handle request create multiple-choice question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.BindJSON(&request); err != nil {
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

func (nova *Nova) HandleCreateQuestionJudgement(c *gin.Context) {
	// create judgement question
	var request QuestionJudgement
	logger.Infof("handle request create judgement question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
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

func (nova *Nova) HandleCreateQuestionEssay(c *gin.Context) {
	// create essay question
	var request QuestionEssay
	logger.Infof("handle request create essay question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.BindJSON(&request); err != nil {
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

func (nova *Nova) HandleDeleteQuestionSingleChoice(c *gin.Context) {
	// delete single choice question
	logger.Infof("handle request delete single-choice question")
	// extract single choice question id from uri
	id := strings.ToLower(c.Param("Id"))
	// request question Id correctness
	logger.Debugf("check single-choice question Id is validate")
	if err := uuid.Validate(id); err != nil {
		nova.response400BadRequest(c, errors.New("single-choice question Id format incorrect"))
		logger.Error("error single-choice question Id is validate")
		return
	}
	logger.Debugf("successfully check single-choice question Id is validate")
	// update data cache by querying single choice question in database
	logger.Debugf("update data cache by querying single-choice questions in database")
	err := nova.querySingleChoiceQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying single-choice questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying single-choice questions in database")
	// check single choice question existence
	logger.Debugf("check single-choice question is validate")
	if !nova.isSingleChoiceQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("single-choice question not found"))
		logger.Error("error check single-choice question is validate")
		return
	}
	logger.Debugf("successfully check single-choice question is validate")
	// delete single choice question from database
	logger.Debugf("delete single-choice question in database")
	if err := nova.deleteSingleChoiceQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Error("error delete single-choice question in database")
		return
	}
	logger.Debugf("successfully delete single-choice question in database")
	// delete single-choice question from data cache
	logger.Debugf("delete single-choice question in data cache")
	nova.deleteSingleChoiceQuestionInDataCache(id)
	// return response
	nova.response204NoContent(c, nil)
	logger.Infof("response status code: %v, body: %v", http.StatusNoContent, nil)
	return
}

func (nova *Nova) HandleDeleteQuestionMultipleChoice(c *gin.Context) {
	// delete multiple-choice question
	logger.Infof("handle request delete multiple-choice question")
	// extract multiple-choice question id from uri
	id := strings.ToLower(c.Param("Id"))
	// request question Id correctness
	logger.Debugf("check multiple-choice question Id is validate")
	if err := uuid.Validate(id); err != nil {
		nova.response400BadRequest(c, errors.New("multiple-choice question Id format incorrect"))
		logger.Error("error multiple-choice question Id is validate")
		return
	}
	logger.Debugf("successfully check multiple-choice question Id is validate")
	// update data cache by querying multiple-choice question in database
	logger.Debugf("update data cache by querying multiple-choice questions in database")
	err := nova.queryMultipleChoiceQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying multiple-choice questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying multiple-choice questions in database")
	// check multiple-choice question existence
	logger.Debugf("check multiple-choice question is validate")
	if !nova.isMultipleChoiceQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("multiple-choice question not found"))
		logger.Error("error check multiple-choice question is validate")
		return
	}
	logger.Debugf("successfully check multiple-choice question is validate")
	// delete multiple-choice question from database
	logger.Debugf("delete multiple-choice question in database")
	if err := nova.deleteMultipleChoiceQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Error("error delete multiple-choice question in database")
		return
	}
	logger.Debugf("successfully delete multiple-choice question in database")
	// delete multiple-choice question from data cache
	logger.Debugf("delete multiple-choice question in data cache")
	nova.deleteMultipleChoiceQuestionInDataCache(id)
	// return response
	nova.response204NoContent(c, nil)
	logger.Infof("response status code: %v, body: %v", http.StatusNoContent, nil)
	return
}

func (nova *Nova) HandleDeleteQuestionJudgement(c *gin.Context) {
	// delete judgement question
	logger.Infof("handle request delete judgement question")
	// extract judgement question id from uri
	id := strings.ToLower(c.Param("Id"))
	// request question Id correctness
	logger.Debugf("check judgement question Id is validate")
	if err := uuid.Validate(id); err != nil {
		nova.response400BadRequest(c, errors.New("judgement question Id format incorrect"))
		logger.Error("error judgement question Id is validate")
		return
	}
	logger.Debugf("successfully check judgement question Id is validate")
	// update data cache by querying judgement question in database
	logger.Debugf("update data cache by querying judgement question in database")
	err := nova.queryJudgementQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying judgement questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying judgement question in database")
	// check judgement question existence
	logger.Debugf("check judgement question is validate")
	if !nova.isJudgementQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("judgement question not found"))
		logger.Error("error check judgement question is validate")
		return
	}
	logger.Debugf("successfully check judgement question is validate")
	// delete judgement question from database
	logger.Debugf("delete judgement question in database")
	if err := nova.deleteJudgementQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Error("error delete judgement question in database")
		return
	}
	logger.Debugf("successfully delete judgement question in database")
	// delete judgement question from data cache
	logger.Debugf("delete judgement question in data cache")
	nova.deleteJudgementQuestionInDataCache(id)
	// return response
	nova.response204NoContent(c, nil)
	logger.Infof("response status code: %v, body: %v", http.StatusNoContent, nil)
	return
}

func (nova *Nova) HandleDeleteQuestionEssay(c *gin.Context) {
	// delete essay question
	logger.Infof("handle request delete essay question")
	// extract essay question id from uri
	id := strings.ToLower(c.Param("Id"))
	// request question Id correctness
	logger.Debugf("check essay question Id is validate")
	if err := uuid.Validate(id); err != nil {
		nova.response400BadRequest(c, errors.New("essay question Id format incorrect"))
		logger.Error("error essay question Id is validate")
		return
	}
	logger.Debugf("successfully check essay question Id is validate")
	// update data cache by querying essay question in database
	logger.Debugf("update data cache by querying essay question in database")
	err := nova.queryEssayQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying essay questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying essay question in database")
	// check essay question existence
	logger.Debugf("check essay question is validate")
	if !nova.isEssayQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("essay question not found"))
		logger.Error("error check essay question is validate")
		return
	}
	logger.Debugf("successfully check essay question is validate")
	// delete essay question from database
	logger.Debugf("delete essay question in database")
	if err := nova.deleteEssayQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Error("error delete essay question in database")
		return
	}
	logger.Debugf("successfully delete essay question in database")
	// delete essay question from data cache
	logger.Debugf("delete essay question in data cache")
	nova.deleteEssayQuestionInDataCache(id)
	// return response
	nova.response204NoContent(c, nil)
	logger.Infof("response status code: %v, body: %v", http.StatusNoContent, nil)
	return
}

func (nova *Nova) HandleModifyQuestionSingleChoice(c *gin.Context) {
	// modify single-choice question
	var request QuestionSingleChoice
	logger.Infof("handle request modify single-choice question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.ShouldBindJSON(&request); err != nil {
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
	if !nova.isSingleChoiceQuestionExisted(strings.ToLower(request.Id)) {
		nova.response404NotFound(c, errors.New("single-choice question not found"))
		logger.Errorf("error check single-choice question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check single-choice question is existed")
	// store modified single-choice question in data cache
	logger.Debugf("store modify single-choice question in data cache")
	response, err := nova.modifySingleChoiceQuestionInDataCache(request)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error store modify single-choice question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully store modify single-choice question in data cache")
	// store modified single-choice question in database
	logger.Debugf("store modify single-choice question in database")
	if err = nova.modifySingleChoiceQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error store modify single-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully store modify single-choice question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleModifyQuestionMultipleChoice(c *gin.Context) {
	// modify multiple-choice question
	var request QuestionMultipleChoice
	logger.Infof("handle request modify multiple-choice question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.ShouldBindJSON(&request); err != nil {
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
	if !nova.isMultipleChoiceQuestionExisted(strings.ToLower(request.Id)) {
		nova.response404NotFound(c, errors.New("multiple-choice question not found"))
		logger.Errorf("error check multiple-choice question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check multiple-choice question is existed")
	// store modified multiple-choice question in data cache
	logger.Debugf("store modify single-choice question in data cache")
	response, err := nova.modifyMultipleChoiceQuestionInDataCache(request)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error store modify multiple-choice question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully store modify multiple-choice question in data cache")
	// store modified multiple-choice question in database
	logger.Debugf("store modify multiple-choice question in database")
	if err = nova.modifyMultipleChoiceQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error store modify multiple-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully store modify multiple-choice question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleModifyQuestionJudgement(c *gin.Context) {
	// modify judgement question
	var request QuestionJudgement
	logger.Infof("handle request modify judgement question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.ShouldBindJSON(&request); err != nil {
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
	if !nova.isJudgementQuestionExisted(strings.ToLower(request.Id)) {
		nova.response404NotFound(c, errors.New("judgement question not found"))
		logger.Errorf("error check judgement question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check judgement question is existed")
	// store modified judgement question in data cache
	logger.Debugf("store modify judgement question in data cache")
	response, err := nova.modifyJudgementQuestionInDataCache(request)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error store modify judgement question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully store modify judgement question in data cache")
	// store modified judgement question in database
	logger.Debugf("store modify judgement question in database")
	if err = nova.modifyJudgementQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error store modify judgement question in database: %v", err)
		return
	}
	logger.Debugf("successfully store modify judgement question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleModifyQuestionEssay(c *gin.Context) {
	// modify essay question
	var request QuestionEssay
	logger.Infof("handle request modify essay question")
	// request body should bind json
	logger.Debugf("request body bind json format")
	if err := c.ShouldBindJSON(&request); err != nil {
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
		logger.Errorf("error update data cache by essay judgement question in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying essay question in database")
	// check essay question existence
	logger.Debugf("check essay question is existed")
	if !nova.isEssayQuestionExisted(strings.ToLower(request.Id)) {
		nova.response404NotFound(c, errors.New("essay question not found"))
		logger.Errorf("error check essay question is existed: %v", err)
		return
	}
	logger.Debugf("successfully check essay question is existed")
	// store modified essay question in data cache
	logger.Debugf("store modify essay question in data cache")
	response, err := nova.modifyEssayQuestionInDataCache(request)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error store modify essay question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully store modify essay question in data cache")
	// store modified essay question in database
	logger.Debugf("store modify essay question in database")
	if err = nova.modifyEssayQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error store modify essay question in database: %v", err)
		return
	}
	logger.Debugf("successfully store modify essay question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleQueryQuestionSingleChoice(c *gin.Context) {
	// query question single-choice
	logger.Infof("handle request query single-choice question")
	// extract questionId from uri
	id := strings.ToLower(c.Param("Id"))
	// request questionId correctness
	logger.Debugf("check questionId is validate")
	if b, _ := nova.isSingleChoiceQuestionValidate(QuestionSingleChoice{Id: id}); !b {
		nova.response400BadRequest(c, errors.New("questionId format incorrect"))
		logger.Errorf("error check questionId is validate")
		return
	}
	logger.Debugf("successfully check questionId is validate")
	// update data cache by querying single-choice questions in database
	logger.Debugf("update data cache by querying single-choice questions in database")
	err := nova.querySingleChoiceQuestionInDatabase(id)
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying single-choice questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying single-choice questions in database")
	// check single-choice question existence
	logger.Debugf("check single-choice question is existed")
	if !nova.isSingleChoiceQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("single-choice question not found"))
		logger.Errorf("error check single-choice question is existed")
		return
	}
	logger.Debugf("successfully check single-choice question is existed")
	// query single-choice question from database
	logger.Debugf("query single-choice question in database")
	if err := nova.querySingleChoiceQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error query single-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully query single-choice question in database")
	// query single-choice question from data cache
	logger.Debugf("query single-choice question in data cache")
	response, err := nova.querySingleChoiceQuestionInDataCache(id)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error query single-choice question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully query single-choice question in data cache")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleQueryQuestionMultipleChoice(c *gin.Context) {
	// query question multiple-choice
	logger.Infof("handle request query multiple-choice question")
	// extract questionId from uri
	id := strings.ToLower(c.Param("Id"))
	// request questionId correctness
	logger.Debugf("check questionId is validate")
	if b, _ := nova.isMultipleChoiceQuestionValidate(QuestionMultipleChoice{Id: id}); !b {
		nova.response400BadRequest(c, errors.New("questionId format incorrect"))
		logger.Errorf("error check questionId is validate")
		return
	}
	logger.Debugf("successfully check questionId is validate")
	// update data cache by querying multiple-choice questions in database
	logger.Debugf("update data cache by querying multiple-choice questions in database")
	err := nova.queryMultipleChoiceQuestionInDatabase(id)
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying multiple-choice questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying multiple-choice questions in database")
	// check multiple-choice question existence
	logger.Debugf("check multiple-choice question is existed")
	if !nova.isMultipleChoiceQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("multiple-choice question not found"))
		logger.Errorf("error check multiple-choice question is existed")
		return
	}
	logger.Debugf("successfully check multiple-choice question is existed")
	// query multiple-choice question from database
	logger.Debugf("query multiple-choice question in database")
	if err := nova.queryMultipleChoiceQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error query multiple-choice question in database: %v", err)
		return
	}
	logger.Debugf("successfully query multiple-choice question in database")
	// query multiple-choice question from data cache
	logger.Debugf("query multiple-choice question in data cache")
	response, err := nova.queryMultipleChoiceQuestionInDataCache(id)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error query multiple-choice question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully query multiple-choice question in data cache")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleQueryQuestionJudgement(c *gin.Context) {
	// query question judgement
	logger.Infof("handle request query judgement question")
	// extract questionId from uri
	id := strings.ToLower(c.Param("Id"))
	// request questionId correctness
	logger.Debugf("check questionId is validate")
	if b, _ := nova.isJudgementQuestionValidate(QuestionJudgement{Id: id}); !b {
		nova.response400BadRequest(c, errors.New("questionId format incorrect"))
		logger.Errorf("error check questionId is validate")
		return
	}
	logger.Debugf("successfully check questionId is validate")
	// update data cache by querying judgement questions in database
	logger.Debugf("update data cache by querying judgement questions in database")
	err := nova.queryJudgementQuestionInDatabase(id)
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying judgement questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying judgement questions in database")
	// check judgement question existence
	logger.Debugf("check judgement question is existed")
	if !nova.isJudgementQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("judgement question not found"))
		logger.Errorf("error check judgement question is existed")
		return
	}
	logger.Debugf("successfully check judgement question is existed")
	// query judgement question from database
	logger.Debugf("query judgement question in database")
	if err := nova.queryJudgementQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error query judgement question in database: %v", err)
		return
	}
	logger.Debugf("successfully query judgement question in database")
	// query judgement question from data cache
	logger.Debugf("query judgement question in data cache")
	response, err := nova.queryJudgementQuestionInDataCache(id)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error query judgement question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully query judgement question in data cache")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleQueryQuestionEssay(c *gin.Context) {
	// query question essay
	logger.Infof("handle request query essay question")
	// extract questionId from uri
	id := strings.ToLower(c.Param("Id"))
	// request questionId correctness
	logger.Debugf("check questionId is validate")
	if b, _ := nova.isEssayQuestionValidate(QuestionEssay{Id: id}); !b {
		nova.response400BadRequest(c, errors.New("questionId format incorrect"))
		logger.Errorf("error check questionId is validate")
		return
	}
	logger.Debugf("successfully check questionId is validate")
	// update data cache by querying essay questions in database
	logger.Debugf("update data cache by querying essay questions in database")
	err := nova.queryEssayQuestionInDatabase(id)
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying essay questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying essay questions in database")
	// check essay question existence
	logger.Debugf("check essay question is existed")
	if !nova.isEssayQuestionExisted(id) {
		nova.response404NotFound(c, errors.New("essay question not found"))
		logger.Errorf("error check essay question is existed")
		return
	}
	logger.Debugf("successfully check essay question is existed")
	// query essay question from database
	logger.Debugf("query essay question in database")
	if err := nova.queryEssayQuestionInDatabase(id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error query essay question in database: %v", err)
		return
	}
	logger.Debugf("successfully query essay question in database")
	// query essay question from data cache
	logger.Debugf("query essay question in data cache")
	response, err := nova.queryEssayQuestionInDataCache(id)
	if err != nil {
		nova.response404NotFound(c, err)
		logger.Errorf("error query essay question in data cache: %v", err)
		return
	}
	logger.Debugf("successfully query essay question in data cache")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleUpdateQuestionSingleChoice(c *gin.Context) {
	// update single-choice question
	var request QuestionSingleChoice
	logger.Infof("handle request update single choice question")
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
	// update data cache by querying single-choice questions in database
	logger.Debugf("update data cache by querying single-choice questions in database")
	err = nova.querySingleChoiceQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying single-choice questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying single-choice questions in database")
	// check single-choice questions existence
	logger.Debugf("check single-choice questions existence")
	if !nova.isSingleChoiceQuestionExisted(strings.ToLower(request.Id)) {
		nova.response403Forbidden(c, errors.New("forbidden replace single-choice question without create it"))
		logger.Errorf("error check single-choice question existence")
		return
	}
	logger.Debugf("successfully check single-choice question existence")
	// store updated single-choice question in data cache
	logger.Debugf("update single-choice question in data cache")
	response := QuestionSingleChoice{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answers:        request.Answers,
		StandardAnswer: request.StandardAnswer,
	}
	if b := nova.updateSingleChoiceQuestionInDataCache(response); !b {
		nova.response404NotFound(c, errors.New("single-choice question not found"))
		logger.Errorf("error update single-choice question in data cache")
		return
	}
	logger.Debugf("successfully update single-choice question in data cache")
	// store update single-choice question in database
	logger.Debugf("update single-choice question in database")
	if err = nova.updateSingleChoiceQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update single-choice question in database")
		return
	}
	logger.Debugf("successfully update single-choice question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleUpdateQuestionMultipleChoice(c *gin.Context) {
	// update multiple-choice question
	var request QuestionMultipleChoice
	logger.Infof("handle request update multiple choice question")
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
	// update data cache by querying multiple-choice questions in database
	logger.Debugf("update data cache by querying multiple-choice questions in database")
	err = nova.queryMultipleChoiceQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying multiple-choice questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying multiple-choice questions in database")
	// check multiple-choice questions existence
	logger.Debugf("check multiple-choice questions existence")
	if !nova.isMultipleChoiceQuestionExisted(strings.ToLower(request.Id)) {
		nova.response403Forbidden(c, errors.New("forbidden replace multiple-choice question without create it"))
		logger.Errorf("error check multiple-choice question existence")
		return
	}
	logger.Debugf("successfully check multiple-choice question existence")
	// store updated multiple-choice question in data cache
	logger.Debugf("update multiple-choice question in data cache")
	response := QuestionMultipleChoice{
		Id:              strings.ToLower(request.Id),
		Title:           request.Title,
		Answers:         request.Answers,
		StandardAnswers: request.StandardAnswers,
	}
	if b := nova.updateMultipleChoiceQuestionInDataCache(response); !b {
		nova.response404NotFound(c, errors.New("multiple-choice question not found"))
		logger.Errorf("error update multiple-choice question in data cache")
		return
	}
	logger.Debugf("successfully update multiple-choice question in data cache")
	// store update multiple-choice question in database
	logger.Debugf("update multiple-choice question in database")
	if err = nova.updateMultipleChoiceQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update multiple-choice question in database")
		return
	}
	logger.Debugf("successfully update multiple-choice question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleUpdateQuestionJudgement(c *gin.Context) {
	// update judgement question
	var request QuestionJudgement
	logger.Infof("handle request update judgement question")
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
	// update data cache by querying judgement in database
	logger.Debugf("update data cache by querying judgement in database")
	err = nova.queryJudgementQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying judgement questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying judgement questions in database")
	// check judgement questions existence
	logger.Debugf("check judgement questions existence")
	if !nova.isJudgementQuestionExisted(strings.ToLower(request.Id)) {
		nova.response403Forbidden(c, errors.New("forbidden replace judgement question without create it"))
		logger.Errorf("error check judgement question existence")
		return
	}
	logger.Debugf("successfully check judgement question existence")
	// store updated judgement question in data cache
	logger.Debugf("update judgement question in data cache")
	response := QuestionJudgement{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answer:         request.Answer,
		StandardAnswer: request.StandardAnswer,
	}
	if b := nova.updateJudgementQuestionInDataCache(response); !b {
		nova.response404NotFound(c, errors.New("judgement question not found"))
		logger.Errorf("error update judgement question in data cache")
		return
	}
	logger.Debugf("successfully update judgement question in data cache")
	// store update judgement question in database
	logger.Debugf("update judgement question in database")
	if err = nova.updateJudgementQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update judgement question in database")
		return
	}
	logger.Debugf("successfully update judgement question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
}

func (nova *Nova) HandleUpdateQuestionEssay(c *gin.Context) {
	// update essay question
	var request QuestionEssay
	logger.Infof("handle request update essay question")
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
	// update data cache by querying essay in database
	logger.Debugf("update data cache by querying essay in database")
	err = nova.queryEssayQuestionsInDatabase()
	if err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update data cache by querying essay questions in database: %v", err)
		return
	}
	logger.Debugf("successfully update data cache by querying essay questions in database")
	// check essay questions existence
	logger.Debugf("check essay questions existence")
	if !nova.isEssayQuestionExisted(strings.ToLower(request.Id)) {
		nova.response403Forbidden(c, errors.New("forbidden replace essay question without create it"))
		logger.Errorf("error check essay question existence")
		return
	}
	logger.Debugf("successfully check essay question existence")
	// store updated judgement question in data cache
	logger.Debugf("update judgement question in data cache")
	response := QuestionEssay{
		Id:             strings.ToLower(request.Id),
		Title:          request.Title,
		Answer:         request.Answer,
		StandardAnswer: request.StandardAnswer,
	}
	if b := nova.updateEssayQuestionInDataCache(response); !b {
		nova.response404NotFound(c, errors.New("essay question not found"))
		logger.Errorf("error update essay question in data cache")
		return
	}
	logger.Debugf("successfully update essay question in data cache")
	// store update essay question in database
	logger.Debugf("update essay question in database")
	if err = nova.updateEssayQuestionInDatabase(response.Id); err != nil {
		nova.response500InternalServerError(c, err)
		logger.Errorf("error update essay question in database")
		return
	}
	logger.Debugf("successfully update essay question in database")
	// return response
	nova.response200OK(c, response)
	logger.Infof("response status code: %v, body: %v", http.StatusOK, response)
	return
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

func (nova *Nova) deleteSingleChoiceQuestionInDataCache(id string) {
	// enable single-choice question cache write lock
	nova.cache.questionsCache.singleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.Unlock()
	// search & delete single-choice question from data cache
	for k, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == id {
			nova.cache.questionsCache.singleChoiceCache.singleChoiceSet = append(nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[:k], nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[k+1:]...)
			break
		}
	}
	return
}

func (nova *Nova) modifySingleChoiceQuestionInDataCache(question QuestionSingleChoice) (QuestionSingleChoice, error) {
	// enable single choice question cache write lock
	nova.cache.questionsCache.singleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.Unlock()
	// replace single choice question in data cache
	for k, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == question.Id {
			nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[k] = question
			return nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[k], nil
		}
	}
	return QuestionSingleChoice{}, errors.New("single choice question not found")
}

func (nova *Nova) querySingleChoiceQuestionInDataCache(id string) (QuestionSingleChoice, error) {
	// enable single-choice question cache read lock
	nova.cache.questionsCache.singleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.RUnlock()
	// search & query single-choice question from data cache
	for k, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == id {
			return nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[k], nil
		}
	}
	return QuestionSingleChoice{}, errors.New("single-choice question not found")
}

func (nova *Nova) updateSingleChoiceQuestionInDataCache(question QuestionSingleChoice) bool {
	// enable single-choice question cache write lock
	nova.cache.questionsCache.singleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.Unlock()
	// replace single-choice question in data cache
	for k, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == question.Id {
			// userName and phoneNumber should not be changed
			nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[k] = question
			return true
		}
	}
	return false
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

func (nova *Nova) deleteSingleChoiceQuestionInDatabase(id string) error {
	// enable single choice question cache read lock
	nova.cache.questionsCache.singleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.RUnlock()
	// search single choice question Id in data cache
	b := false
	for _, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == id {
			b = true
			break
		}
	}
	if !b {
		return errors.New("single-choice question not found")
	}
	// delete single-choice question in database
	if err := nova.db.DeleteQuestionSingleChoice(id); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) modifySingleChoiceQuestionInDatabase(id string) error {
	// enable single choice question cache read lock
	nova.cache.questionsCache.singleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.RUnlock()
	// search single choice question Id in data cache
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
	// update single choice question in database
	if err := nova.db.UpdateQuestionSingleChoice(&question); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) querySingleChoiceQuestionInDatabase(id string) error {
	// enable single-choice question cache write lock
	nova.cache.questionsCache.singleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.Unlock()
	// query single-choice question from database
	question, err := nova.db.QueryQuestionSingleChoice(id)
	if err != nil {
		return err
	}
	// update single-choice question in data cache
	for k, v := range nova.cache.questionsCache.singleChoiceCache.singleChoiceSet {
		if v.Id == id {
			nova.cache.questionsCache.singleChoiceCache.singleChoiceSet[k] = *question
			break
		}
	}
	return nil
}

func (nova *Nova) updateSingleChoiceQuestionInDatabase(id string) error {
	// enable single-choice question cache read lock
	nova.cache.questionsCache.singleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.singleChoiceCache.mutex.RUnlock()
	// search single-choice question in data cache
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
	// update single-choice question in database
	if err := nova.db.UpdateQuestionSingleChoice(&question); err != nil {
		return err
	}
	return nil
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

func (nova *Nova) createMultipleChoiceQuestionInDataCache(question QuestionMultipleChoice) {
	// enable multiple-choice question cache write lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.Unlock()
	// append multiple-choice question in data cache
	nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet = append(nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet, question)
	return
}

func (nova *Nova) deleteMultipleChoiceQuestionInDataCache(id string) {
	// enable multiple-choice question cache write lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.Unlock()
	// search & delete multiple-choice question from data cache
	for k, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == id {
			nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet = append(nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[:k], nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[k+1:]...)
			break
		}
	}
	return
}

func (nova *Nova) modifyMultipleChoiceQuestionInDataCache(question QuestionMultipleChoice) (QuestionMultipleChoice, error) {
	// enable multiple choice question cache write lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.Unlock()
	// replace multiple choice question in data cache
	for k, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == question.Id {
			nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[k] = question
			return nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[k], nil
		}
	}
	return QuestionMultipleChoice{}, errors.New("multiple choice question not found")
}

func (nova *Nova) queryMultipleChoiceQuestionInDataCache(id string) (QuestionMultipleChoice, error) {
	// enable multiple-choice question cache read lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.RUnlock()
	// search & query multiple-choice question from data cache
	for k, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == id {
			return nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[k], nil
		}
	}
	return QuestionMultipleChoice{}, errors.New("multiple-choice question not found")
}

func (nova *Nova) updateMultipleChoiceQuestionInDataCache(question QuestionMultipleChoice) bool {
	// enable multiple-choice question cache write lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.Unlock()
	// replace multiple-choice question in data cache
	for k, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == question.Id {
			// userName and phoneNumber should not be changed
			nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[k] = question
			return true
		}
	}
	return false
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

func (nova *Nova) deleteMultipleChoiceQuestionInDatabase(id string) error {
	// enable multiple-choice question cache read lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.RUnlock()
	// search multiple-choice question Id in data cache
	b := false
	for _, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == id {
			b = true
			break
		}
	}
	if !b {
		return errors.New("multiple-choice question not found")
	}
	// delete multiple-choice question in database
	if err := nova.db.DeleteQuestionMultipleChoice(id); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) modifyMultipleChoiceQuestionInDatabase(id string) error {
	// enable multiple choice question cache read lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.RUnlock()
	// search multiple choice question Id in data cache
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
	// update multiple choice question in database
	if err := nova.db.UpdateQuestionMultipleChoice(&question); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) queryMultipleChoiceQuestionInDatabase(id string) error {
	// enable multiple-choice question cache write lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.Lock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.Unlock()
	// query multiple-choice question from database
	question, err := nova.db.QueryQuestionMultipleChoice(id)
	if err != nil {
		return err
	}
	// update multiple-choice question in data cache
	for k, v := range nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet {
		if v.Id == id {
			nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet[k] = *question
			break
		}
	}
	return nil
}

func (nova *Nova) updateMultipleChoiceQuestionInDatabase(id string) error {
	// enable multiple-choice question cache read lock
	nova.cache.questionsCache.multipleChoiceCache.mutex.RLock()
	defer nova.cache.questionsCache.multipleChoiceCache.mutex.RUnlock()
	// search multiple-choice question in data cache
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
	// update multiple-choice question in database
	if err := nova.db.UpdateQuestionMultipleChoice(&question); err != nil {
		return err
	}
	return nil
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

func (nova *Nova) createJudgementQuestionInDataCache(question QuestionJudgement) {
	// enable judgement question cache write lock
	nova.cache.questionsCache.judgementCache.mutex.Lock()
	defer nova.cache.questionsCache.judgementCache.mutex.Unlock()
	// append judgement question in data cache
	nova.cache.questionsCache.judgementCache.judgementSet = append(nova.cache.questionsCache.judgementCache.judgementSet, question)
	return
}

func (nova *Nova) deleteJudgementQuestionInDataCache(id string) {
	// enable judgement question cache write lock
	nova.cache.questionsCache.judgementCache.mutex.Lock()
	defer nova.cache.questionsCache.judgementCache.mutex.Unlock()
	// search & delete judgement question from data cache
	for k, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == id {
			nova.cache.questionsCache.judgementCache.judgementSet = append(nova.cache.questionsCache.judgementCache.judgementSet[:k], nova.cache.questionsCache.judgementCache.judgementSet[k+1:]...)
			break
		}
	}
	return
}

func (nova *Nova) modifyJudgementQuestionInDataCache(question QuestionJudgement) (QuestionJudgement, error) {
	// enable judgement question cache write lock
	nova.cache.questionsCache.judgementCache.mutex.Lock()
	defer nova.cache.questionsCache.judgementCache.mutex.Unlock()
	// replace judgement question in data cache
	for k, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == question.Id {
			nova.cache.questionsCache.judgementCache.judgementSet[k] = question
			return nova.cache.questionsCache.judgementCache.judgementSet[k], nil
		}
	}
	return QuestionJudgement{}, errors.New("judgement question not found")
}

func (nova *Nova) queryJudgementQuestionInDataCache(id string) (QuestionJudgement, error) {
	// enable judgement question cache read lock
	nova.cache.questionsCache.judgementCache.mutex.RLock()
	defer nova.cache.questionsCache.judgementCache.mutex.RUnlock()
	// search & query judgement question from data cache
	for k, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == id {
			return nova.cache.questionsCache.judgementCache.judgementSet[k], nil
		}
	}
	return QuestionJudgement{}, errors.New("judgement question not found")
}

func (nova *Nova) updateJudgementQuestionInDataCache(question QuestionJudgement) bool {
	// enable judgement question cache write lock
	nova.cache.questionsCache.judgementCache.mutex.Lock()
	defer nova.cache.questionsCache.judgementCache.mutex.Unlock()
	// replace judgement question in data cache
	for k, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == question.Id {
			// userName and phoneNumber should not be changed
			nova.cache.questionsCache.judgementCache.judgementSet[k] = question
			return true
		}
	}
	return false
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

func (nova *Nova) deleteJudgementQuestionInDatabase(id string) error {
	// enable judgement question cache read lock
	nova.cache.questionsCache.judgementCache.mutex.RLock()
	defer nova.cache.questionsCache.judgementCache.mutex.RUnlock()
	// search judgement question Id in data cache
	b := false
	for _, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == id {
			b = true
			break
		}
	}
	if !b {
		return errors.New("judgement question not found")
	}
	// delete judgement question in database
	if err := nova.db.DeleteQuestionJudgement(id); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) modifyJudgementQuestionInDatabase(id string) error {
	// enable judgement question cache read lock
	nova.cache.questionsCache.judgementCache.mutex.RLock()
	defer nova.cache.questionsCache.judgementCache.mutex.RUnlock()
	// search judgement question Id in data cache
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
	// update judgement question in database
	if err := nova.db.UpdateQuestionJudgement(&question); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) queryJudgementQuestionInDatabase(id string) error {
	// enable judgement question cache write lock
	nova.cache.questionsCache.judgementCache.mutex.Lock()
	defer nova.cache.questionsCache.judgementCache.mutex.Unlock()
	// query judgement question from database
	question, err := nova.db.QueryQuestionJudgement(id)
	if err != nil {
		return err
	}
	// update judgement question in data cache
	for k, v := range nova.cache.questionsCache.judgementCache.judgementSet {
		if v.Id == id {
			nova.cache.questionsCache.judgementCache.judgementSet[k] = *question
			break
		}
	}
	return nil
}

func (nova *Nova) updateJudgementQuestionInDatabase(id string) error {
	// enable judgement question cache read lock
	nova.cache.questionsCache.judgementCache.mutex.RLock()
	defer nova.cache.questionsCache.judgementCache.mutex.RUnlock()
	// search judgement question in data cache
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
	// update judgement question in database
	if err := nova.db.UpdateQuestionJudgement(&question); err != nil {
		return err
	}
	return nil
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

func (nova *Nova) createEssayQuestionInDataCache(question QuestionEssay) {
	// enable essay question cache write lock
	nova.cache.questionsCache.essayCache.mutex.Lock()
	defer nova.cache.questionsCache.essayCache.mutex.Unlock()
	// append essay question in data cache
	nova.cache.questionsCache.essayCache.essaySet = append(nova.cache.questionsCache.essayCache.essaySet, question)
	return
}

func (nova *Nova) deleteEssayQuestionInDataCache(id string) {
	// enable essay question cache write lock
	nova.cache.questionsCache.essayCache.mutex.Lock()
	defer nova.cache.questionsCache.essayCache.mutex.Unlock()
	// search & delete essay question from data cache
	for k, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == id {
			nova.cache.questionsCache.essayCache.essaySet = append(nova.cache.questionsCache.essayCache.essaySet[:k], nova.cache.questionsCache.essayCache.essaySet[k+1:]...)
			break
		}
	}
	return
}

func (nova *Nova) modifyEssayQuestionInDataCache(question QuestionEssay) (QuestionEssay, error) {
	// enable essay question cache write lock
	nova.cache.questionsCache.essayCache.mutex.Lock()
	defer nova.cache.questionsCache.essayCache.mutex.Unlock()
	// replace essay question in data cache
	for k, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == question.Id {
			nova.cache.questionsCache.essayCache.essaySet[k] = question
			return nova.cache.questionsCache.essayCache.essaySet[k], nil
		}
	}
	return QuestionEssay{}, errors.New("essay question not found")
}

func (nova *Nova) queryEssayQuestionInDataCache(id string) (QuestionEssay, error) {
	// enable essay question cache read lock
	nova.cache.questionsCache.essayCache.mutex.RLock()
	defer nova.cache.questionsCache.essayCache.mutex.RUnlock()
	// search & query essay question from data cache
	for k, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == id {
			return nova.cache.questionsCache.essayCache.essaySet[k], nil
		}
	}
	return QuestionEssay{}, errors.New("essay question not found")
}

func (nova *Nova) updateEssayQuestionInDataCache(question QuestionEssay) bool {
	// enable essay question cache write lock
	nova.cache.questionsCache.essayCache.mutex.Lock()
	defer nova.cache.questionsCache.essayCache.mutex.Unlock()
	// replace essay question in data cache
	for k, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == question.Id {
			// userName and phoneNumber should not be changed
			nova.cache.questionsCache.essayCache.essaySet[k] = question
			return true
		}
	}
	return false
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

func (nova *Nova) deleteEssayQuestionInDatabase(id string) error {
	// enable essay question cache read lock
	nova.cache.questionsCache.essayCache.mutex.RLock()
	defer nova.cache.questionsCache.essayCache.mutex.RUnlock()
	// search essay question Id in data cache
	b := false
	for _, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == id {
			b = true
			break
		}
	}
	if !b {
		return errors.New("essay question not found")
	}
	// delete essay question in database
	if err := nova.db.DeleteQuestionEssay(id); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) modifyEssayQuestionInDatabase(id string) error {
	// enable essay question cache read lock
	nova.cache.questionsCache.essayCache.mutex.RLock()
	defer nova.cache.questionsCache.essayCache.mutex.RUnlock()
	// search essay question Id in data cache
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
	// update essay question in database
	if err := nova.db.UpdateQuestionEssay(&question); err != nil {
		return err
	}
	return nil
}

func (nova *Nova) queryEssayQuestionInDatabase(id string) error {
	// enable essay question cache write lock
	nova.cache.questionsCache.essayCache.mutex.Lock()
	defer nova.cache.questionsCache.essayCache.mutex.Unlock()
	// query essay question from database
	question, err := nova.db.QueryQuestionEssay(id)
	if err != nil {
		return err
	}
	// update essay question in data cache
	for k, v := range nova.cache.questionsCache.essayCache.essaySet {
		if v.Id == id {
			nova.cache.questionsCache.essayCache.essaySet[k] = *question
			break
		}
	}
	return nil
}

func (nova *Nova) updateEssayQuestionInDatabase(id string) error {
	// enable essay question cache read lock
	nova.cache.questionsCache.essayCache.mutex.RLock()
	defer nova.cache.questionsCache.essayCache.mutex.RUnlock()
	// search essay question in data cache
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
	// update essay question in database
	if err := nova.db.UpdateQuestionEssay(&question); err != nil {
		return err
	}
	return nil
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

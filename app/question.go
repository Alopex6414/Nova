package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nova/logger"
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

}

func (nova *Nova) HandleDeleteQuestion(c *gin.Context) {

}

func (nova *Nova) HandleModifyQuestion(c *gin.Context) {

}

func (nova *Nova) HandleQueryQuestion(c *gin.Context) {

}

func (nova *Nova) HandleUpdateQuestion(c *gin.Context) {

}

package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (nova *Nova) response200OK(c *gin.Context, body any) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, body)
	return
}

func (nova *Nova) response201Created(c *gin.Context, body any) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, body)
	return
}

func (nova *Nova) response204NoContent(c *gin.Context, body any) {
	c.Status(http.StatusNoContent)
	return
}

func (nova *Nova) response400BadRequest(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Bad Request"
	problemDetails.Type = "Client Error"
	problemDetails.Status = http.StatusBadRequest
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusBadRequest, problemDetails)
	return
}

func (nova *Nova) response403Forbidden(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Forbidden"
	problemDetails.Type = "Client Error"
	problemDetails.Status = http.StatusForbidden
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusForbidden, problemDetails)
	return
}

func (nova *Nova) response404NotFound(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Not Found"
	problemDetails.Type = "Client Error"
	problemDetails.Status = http.StatusNotFound
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusNotFound, problemDetails)
	return
}

func (nova *Nova) response409Conflict(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Conflict"
	problemDetails.Type = "Client Error"
	problemDetails.Status = http.StatusConflict
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusConflict, problemDetails)
	return
}

func (nova *Nova) response412PreconditionFailed(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "PreconditionFailed"
	problemDetails.Type = "Client Error"
	problemDetails.Status = http.StatusPreconditionFailed
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusPreconditionFailed, problemDetails)
	return
}

func (nova *Nova) response417ExpectationFailed(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "ExpectationFailed"
	problemDetails.Type = "Client Error"
	problemDetails.Status = http.StatusExpectationFailed
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusExpectationFailed, problemDetails)
	return
}

func (nova *Nova) response500InternalServerError(c *gin.Context, err error) {
	var problemDetails ProblemDetails
	problemDetails.Title = "Internal Server Error"
	problemDetails.Type = "Server Error"
	problemDetails.Status = http.StatusInternalServerError
	problemDetails.Cause = err.Error()
	c.Header("Content-Type", "application/problem+json")
	c.JSON(http.StatusInternalServerError, problemDetails)
	return
}

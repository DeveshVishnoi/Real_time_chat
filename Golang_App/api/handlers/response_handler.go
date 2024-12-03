package handlers

import (
	"realtime_chat/models"

	"github.com/gin-gonic/gin"
)

// ReturnResponse will be used as Response template to send the response for API
func ConstructResponse(statusCode int, message string, data any, err error) *models.Response {
	var error string
	if err != nil {
		error = err.Error()
	}
	return &models.Response{StatusCode: statusCode, Message: message, Data: data, Error: error}
}

func RespondError(ctx *gin.Context, statusCode int, message string, err error) {
	ctx.JSON(statusCode, models.Response{
		StatusCode: statusCode,
		Message:    message,
		Error:      err.Error(),
	})
}

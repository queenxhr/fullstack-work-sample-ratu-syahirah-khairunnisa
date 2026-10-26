package utils

import (
	"net/http"

	appErrors "fullstack-backend/internal/errors"
	"github.com/gin-gonic/gin"
)

// DefaultResponse represents standard JSON response
type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseError represents standard error response
type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func JSON(ctx *gin.Context, statusCode int, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	response := DefaultResponse{
		Status:  http.StatusText(statusCode),
		Message: msg,
		Data:    data,
	}

	ctx.JSON(statusCode, response)
}

func SendJSONErrorResponse(ctx *gin.Context, err error) {
	var statusCode = http.StatusInternalServerError
	var httpStatus = http.StatusText(http.StatusInternalServerError)
	var message = "Internal Server Error"

	if err != nil {
		message = err.Error()
	}

	if dynamicErr, ok := err.(*appErrors.DynamicError); ok {
		statusCode = dynamicErr.StatusCode
		httpStatus = http.StatusText(statusCode)
		message = dynamicErr.Message
	}

	response := ResponseError{
		Status:  httpStatus,
		Message: message,
	}

	ctx.AbortWithStatusJSON(statusCode, response)
}

package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	ErrBadRequest ErrorCode = http.StatusBadRequest
	ErrNotFound   ErrorCode = http.StatusNotFound
	ErrConflict   ErrorCode = http.StatusConflict
	ErrUnauthorized ErrorCode = http.StatusUnauthorized
	ErrInternal   ErrorCode = http.StatusInternalServerError
)

type DynamicError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *DynamicError) Error() string {
	return e.Message
}

func New(code ErrorCode, message string) error {
	if message == "" {
		message = http.StatusText(int(code))
	}
	return &DynamicError{
		StatusCode: int(code),
		Message:    message,
	}
}

func Newf(code ErrorCode, format string, args ...interface{}) error {
	return &DynamicError{
		StatusCode: int(code),
		Message:    fmt.Sprintf(format, args...),
	}
}

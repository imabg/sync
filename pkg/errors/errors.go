package errors

import (
	"fmt"
	"net/http"
)

const (
	VALIDATION_ERROR = "validation error"
	DATABASE_ERROR = "database error"
	NOT_FOUND = "not found"
	INTERNAL_SERVER_ERROR = "internal server error"
	CONFLICT_ERROR = "conflict error"
)

type myError struct {
	code int
	msg  string
}

type CustomError struct {
	Code int `json:"code"`
	Msg string `json:"message"`
	AdditionalCode string `json:"additionalCode,omitempty"`
	AdditionalMsg string `json:"additionalMessage,omitempty"`
}



func NewCustomError(code int, msg string, additionalCode string, additionalMsg string) *CustomError {
	return &CustomError{
		Code: code,
		Msg: msg,
		AdditionalCode: additionalCode,
		AdditionalMsg: additionalMsg,
	}
}

func (e CustomError) Error() string {
	return fmt.Sprint(e.Msg)
}

func BadRequestError(msg string) *CustomError{
	return &CustomError{
		Code: http.StatusBadRequest,
		Msg: msg,
	}
}

func ConflictError(msg string) *CustomError {
	return &CustomError{
		Code: http.StatusConflict,
		Msg: msg,
	}
}

func NotFound(msg string) *CustomError {
	return &CustomError{
		Code: http.StatusNotFound,
		Msg: msg,
	}
}

func AuthorizationError(msg string) *CustomError {
	return &CustomError{
		Code: http.StatusUnauthorized,
		Msg: msg,
	}
}
package errors

import "fmt"

const (
	VALIDATION_ERROR = "validation error"
	DATABASE_ERROR = "database error"
	NOT_FOUND = "not found"
	INTERNAL_SERVER_ERROR = "internal server error"
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
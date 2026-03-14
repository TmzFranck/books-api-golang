package utils

import "net/http"

var (
	ErrBadRequest          = NewError(http.StatusBadRequest)
	ErrForbidden           = NewError(http.StatusForbidden)
	ErrNotFound            = NewError(http.StatusNotFound)
	ErrInternalServerError = NewError(http.StatusInternalServerError)
	ErrUnauthorized        = NewError(http.StatusUnauthorized)
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int) *Error {
	return &Error{
		Code:    code,
		Message: http.StatusText(code),
	}
}

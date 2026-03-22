package utils

import "net/http"

var (
	ErrBadRequest          = NewError(http.StatusBadRequest)
	ErrForbidden           = NewError(http.StatusForbidden)
	ErrNotFound            = NewError(http.StatusNotFound)
	ErrInternalServerError = NewError(http.StatusInternalServerError)
	ErrUnauthorized        = NewError(http.StatusUnauthorized)
)

// Error represents an HTTP error with a message and status code
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error returns the error message
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new Error with the given status code and message
func NewError(code int) *Error {
	return &Error{
		Code:    code,
		Message: http.StatusText(code),
	}
}

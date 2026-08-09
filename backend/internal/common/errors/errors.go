package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.StatusCode, e.Message)
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

var (
	ErrBadRequest   = NewAPIError(http.StatusBadRequest, "Bad Request")
	ErrUnauthorized = NewAPIError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden    = NewAPIError(http.StatusForbidden, "Forbidden")
	ErrNotFound     = NewAPIError(http.StatusNotFound, "Not Found")
	ErrConflict     = NewAPIError(http.StatusConflict, "Conflict")
	ErrInternal     = NewAPIError(http.StatusInternalServerError, "Internal Server Error")
)

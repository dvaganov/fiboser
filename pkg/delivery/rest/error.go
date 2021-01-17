package rest

import "net/http"

type (
	httpError struct {
		Status     int    `json:"-"`
		Message    string `json:"message"`
		innerError error
	}
)

func NewRequestError() httpError {
	return httpError{http.StatusBadRequest, "bad request", nil}
}

func NewInternalError(err error) httpError {
	return httpError{http.StatusInternalServerError, "internal server error", err}
}

func (e httpError) Error() string {
	return e.Message
}

func (e httpError) Unwrap() error {
	return e.innerError
}

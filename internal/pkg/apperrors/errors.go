package apperrors

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidRequestParam = errors.New("invalid request param")
	ErrInvalidRequestBody  = errors.New("invalid request body")
)

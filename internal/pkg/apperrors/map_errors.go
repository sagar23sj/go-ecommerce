package apperrors

import "net/http"

func MapError(err error) (statusCode int, errResponse error) {
	if _, ok := err.(ProductNotFound); ok {
		return http.StatusUnprocessableEntity, err
	}

	return
}

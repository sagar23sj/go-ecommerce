package apperrors

import "net/http"

func MapError(err error) (statusCode int, errResponse error) {
	if _, ok := err.(ProductNotFound); ok {
		return http.StatusBadRequest, err
	}

	if _, ok := err.(ProductQuantityInsufficient); ok {
		return http.StatusUnprocessableEntity, err
	}

	if _, ok := err.(ProductQuantityExceeded); ok {
		return http.StatusUnprocessableEntity, err
	}

	if _, ok := err.(OrderStatusInvalid); ok {
		return http.StatusUnprocessableEntity, err
	}

	if _, ok := err.(OrderUpdationInvalid); ok {
		return http.StatusUnprocessableEntity, err
	}

	if _, ok := err.(OrderNotFound); ok {
		return http.StatusBadRequest, err
	}

	return http.StatusInternalServerError, err
}

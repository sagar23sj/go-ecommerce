package apperrors

import "net/http"

func MapError(err error) (statusCode int, errResponse error) {
	switch err.(type) {
	case ProductNotFound:
		return http.StatusBadRequest, err
	case ProductQuantityInsufficient:
		return http.StatusUnprocessableEntity, err
	case ProductQuantityExceeded:
		return http.StatusUnprocessableEntity, err
	case ProductQuantityInvalid:
		return http.StatusUnprocessableEntity, err
	case OrderStatusInvalid:
		return http.StatusUnprocessableEntity, err
	case OrderUpdationInvalid:
		return http.StatusUnprocessableEntity, err
	case OrderNotFound:
		return http.StatusBadRequest, err
	default:
		return http.StatusInternalServerError, err
	}
}

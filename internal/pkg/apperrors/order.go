package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrNoProductsToOrder = errors.New("no products to order")
)

type OrderNotFound struct {
	ID int64
}

func (o OrderNotFound) Error() string {
	return fmt.Sprintf("order not found with id: %d", o.ID)
}

type OrderStatusInvalid struct {
	ID int64
}

func (o OrderStatusInvalid) Error() string {
	return fmt.Sprintf("invalid status for order with id: %d", o.ID)
}

type OrderUpdationInvalid struct {
	ID             int64
	CurrentState   string
	RequestedState string
}

func (o OrderUpdationInvalid) Error() string {
	return fmt.Sprintf("order updation invalid for order with id: %d, current_state: %s, requested_state: %s", o.ID, o.CurrentState, o.RequestedState)
}

package apperrors

import "fmt"

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
	return fmt.Sprintf("invalid status for order with: %d", o.ID)
}

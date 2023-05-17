package apperrors

import "fmt"

type ProductNotFound struct {
	ID int64
}

func (p ProductNotFound) Error() string {
	return fmt.Sprintf("product not found with id: %d", p.ID)
}

type ProductQuantityInsufficient struct {
	ID                int64
	QuantityAsked     int64
	QuantityRemaining int64
}

func (p ProductQuantityInsufficient) Error() string {
	return fmt.Sprintf("product quantity insufficient for id: %d, quantity_remaining : %d and quantity_asked : %d", p.ID, p.QuantityRemaining, p.QuantityAsked)
}

type ProductQuantityExceeded struct {
	ID            int64
	QuantityLimit int64
	QuantityAsked int64
}

func (p ProductQuantityExceeded) Error() string {
	return fmt.Sprintf("product quantity exceeded for id: %d, quantity_limit : %d and quantity_asked : %d", p.ID, p.QuantityLimit, p.QuantityAsked)
}

type ProductQuantityInvalid struct {
	ID       int64
	Quantity int64
}

func (p ProductQuantityInvalid) Error() string {
	return fmt.Sprintf("product quantity invalid for id: %d, quantity : %d", p.ID, p.Quantity)
}

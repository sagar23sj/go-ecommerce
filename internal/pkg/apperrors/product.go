package apperrors

import "fmt"

type ProductNotFound struct {
	ID int64
}

func (p ProductNotFound) Error() string {
	return fmt.Sprintf("product not found with id: %d", p.ID)
}

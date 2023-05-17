package dto

import (
	"fmt"
	"time"
)

type Order struct {
	ID                 int64         `json:"id"`
	Products           []ProductInfo `json:"products,omitempty"`
	Amount             float64       `json:"amount"`
	DiscountPercentage float64       `json:"discount_percent"`
	FinalAmount        float64       `json:"final_amount"`
	Status             string        `json:"status"`
	DispatchedAt       time.Time     `json:"dispatched_at,omitempty"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

type ProductInfo struct {
	ProductID int64 `json:"product_id,omitempty"`
	Quantity  int64 `json:"quantity,omitempty"`
}

type CreateOrderRequest struct {
	Products []ProductInfo `json:"products"`
}

type UpdateOrderStatusRequest struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

func (req *CreateOrderRequest) Validate() error {
	for _, p := range req.Products {
		if p.Quantity <= 0 {
			return fmt.Errorf("invalid request, product quantity negative for product_id : %d", p.ProductID)
		}
	}

	return nil
}

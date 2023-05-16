package dto

import (
	"fmt"
	"time"
)

type Order struct {
	ID                 int64         `json:"id,omitempty"`
	Products           []ProductInfo `json:"products,omitempty"`
	Amount             float64       `json:"amount,omitempty"`
	DiscountPercentage float64       `json:"discount_percent,omitempty"`
	DiscountedAmount   float64       `json:"discounted_amount,omitempty"`
	Status             string        `json:"status,omitempty"`
	CreatedAt          time.Time     `json:"created_at,omitempty"`
	UpdatedAt          time.Time     `json:"updated_at,omitempty"`
}

type ProductInfo struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type CreateOrderRequest struct {
	Products []ProductInfo `json:"products"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (req *CreateOrderRequest) Validate() error {
	for _, p := range req.Products {
		if p.Quantity <= 0 {
			return fmt.Errorf("invalid request, product quantity negative for product_id : %d", p.ProductID)
		}
	}

	return nil
}

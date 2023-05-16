package dto

import "time"

type Order struct {
	ID                 int64         `json:"id"`
	Products           []ProductInfo `json:"products"`
	Amount             float64       `json:"amount"`
	DiscountPercentage int64         `json:"discount_percent"`
	DiscountedAmount   float64       `json:"discounted_amount"`
	Status             string        `json:"status"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
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

package dto

import (
	"errors"
	"fmt"
	"time"

	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
)

type Order struct {
	ID                 int64         `json:"id"`
	Products           []ProductInfo `json:"products,omitempty"`
	Amount             float64       `json:"amount"`
	DiscountPercentage float64       `json:"discount_percent"`
	FinalAmount        float64       `json:"final_amount"`
	Status             string        `json:"status"`
	DispatchedAt       *time.Time    `json:"dispatched_at,omitempty"`
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

	if len(req.Products) <= 0 {
		return apperrors.ErrNoProductsToOrder
	}

	//map[ProductID]bool
	productMap := make(map[int64]bool)
	for _, p := range req.Products {
		if _, ok := productMap[p.ProductID]; ok {
			return fmt.Errorf("invalid request, duplicate product found with product_id : %d", p.ProductID)
		}

		if p.Quantity <= 0 {
			return fmt.Errorf("invalid request, product quantity negative for product_id : %d", p.ProductID)
		}

		productMap[p.ProductID] = true
	}

	return nil
}

func (req *UpdateOrderStatusRequest) Validate() error {
	if req.OrderID == 0 {
		return errors.New("order_id cannot be empty")
	}

	if req.Status == "" {
		return errors.New("status cannot be empty")
	}

	return nil
}

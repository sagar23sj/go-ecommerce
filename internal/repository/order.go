package repository

import (
	"context"
	"time"
)

type OrderStorer interface {
	RepositoryTransaction

	GetOrderByID(ctx context.Context, tx Transaction, orderID int64) (Order, error)
	CreateOrder(ctx context.Context, tx Transaction, order Order) (Order, error)
	UpdateOrderStatus(ctx context.Context, tx Transaction, orderID int64, status string) error
	UpdateOrderDispatchDate(ctx context.Context, tx Transaction, orderID int64, dispatchedAt time.Time) error
	ListOrders(ctx context.Context, tx Transaction) ([]Order, error)
}

type Order struct {
	ID                 uint `storm:"id,increment"`
	Amount             float64
	DiscountPercentage float64
	FinalAmount        float64
	Status             string
	DispatchedAt       time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

package repository

import (
	"context"
	"time"
)

type OrderItemStorer interface {
	RepositoryTransaction

	GetOrderItemsByOrderID(ctx context.Context, tx Transaction, orderID int64) ([]OrderItem, error)
	StoreOrderItems(ctx context.Context, tx Transaction, orderItems []OrderItem) error
}

type OrderItem struct {
	ID        uint `storm:"id,increment"`
	OrderID   int64
	ProductID int64
	Quantity  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

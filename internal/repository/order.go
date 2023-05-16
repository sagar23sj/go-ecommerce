package repository

import (
	"context"

	"gorm.io/gorm"
)

type orderStore struct {
	BaseRepository
	OrderDetailsRepo OrderDetailsStorer
}

type OrderStorer interface {
	GetOrderByID(ctx context.Context, orderID int64) (Order, error)
	CreateOrder(ctx context.Context, order Order) (OrderDetails, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) (Order, error)
}

func NewOrderRepo(db *gorm.DB) OrderStorer {
	return &orderStore{
		BaseRepository:   BaseRepository{db},
		OrderDetailsRepo: NewOrderDetailsRepo(db),
	}
}

type Order struct {
	gorm.Model
	Amount             float64
	DiscountPercentage int64
	DiscountedAmount   float64
	Status             string
}

func (os *orderStore) GetOrderByID(ctx context.Context, orderID int64) (Order, error) {
	return Order{}, nil
}

func (os *orderStore) CreateOrder(ctx context.Context, order Order) (OrderDetails, error) {
	return OrderDetails{}, nil
}

func (os *orderStore) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (Order, error) {
	return Order{}, nil
}

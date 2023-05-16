package repository

import (
	"context"

	"gorm.io/gorm"
)

type orderDetailsStore struct {
	BaseRepository
}

type OrderDetailsStorer interface {
	GetOrderDetailsByOrderID(ctx context.Context, orderID int64) (OrderDetails, error)
	StoreOrderDetails(ctx context.Context, orderDetails OrderDetails) error
}

func NewOrderDetailsRepo(db *gorm.DB) OrderDetailsStorer {
	return &orderDetailsStore{
		BaseRepository: BaseRepository{db},
	}
}

type OrderDetails struct {
	gorm.Model
	OrderID   int64
	ProductID int64
	Quantity  int64
	Category  string
	Price     float64
}

func (ods *orderDetailsStore) GetOrderDetailsByOrderID(ctx context.Context, orderID int64) (OrderDetails, error) {
	return OrderDetails{}, nil
}

func (ods *orderDetailsStore) StoreOrderDetails(ctx context.Context, orderDetails OrderDetails) error {
	return nil
}

package repository

import (
	"context"

	"gorm.io/gorm"
)

type orderItemStore struct {
	BaseRepository
}

type OrderItemStorer interface {
	RepositoryTransaction

	GetOrderItemsByOrderID(ctx context.Context, tx *gorm.DB, orderID int64) ([]OrderItem, error)
	StoreOrderItems(ctx context.Context, tx *gorm.DB, orderItems []OrderItem) error
}

func NewOrderItemRepo(db *gorm.DB) OrderItemStorer {
	return &orderItemStore{
		BaseRepository: BaseRepository{db},
	}
}

type OrderItem struct {
	gorm.Model
	OrderID   int64
	ProductID int64
	Quantity  int64
}

func (ods *orderItemStore) GetOrderItemsByOrderID(ctx context.Context, tx *gorm.DB, orderID int64) ([]OrderItem, error) {
	orderItemList := make([]OrderItem, 0)

	queryExecutor := ods.initiateQueryExecutor(tx)
	err := queryExecutor.Find(&orderItemList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return orderItemList, err
	}

	return orderItemList, nil
}

func (ods *orderItemStore) StoreOrderItems(ctx context.Context, tx *gorm.DB, orderItems []OrderItem) error {
	queryExecutor := ods.initiateQueryExecutor(tx)
	err := queryExecutor.Create(&orderItems).Error
	if err != nil {
		return err
	}

	return nil
}

package repository

import (
	"context"
	"time"

	"github.com/asdine/storm/v3"
)

type orderItemStore struct {
	BaseRepository
}

type OrderItemStorer interface {
	RepositoryTransaction

	GetOrderItemsByOrderID(ctx context.Context, tx Transaction, orderID int64) ([]OrderItem, error)
	StoreOrderItems(ctx context.Context, tx Transaction, orderItems []OrderItem) error
}

func NewOrderItemRepo(db *storm.DB) OrderItemStorer {
	return &orderItemStore{
		BaseRepository: BaseRepository{db},
	}
}

type OrderItem struct {
	ID        uint `storm:"id,increment"`
	OrderID   int64
	ProductID int64
	Quantity  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ods *orderItemStore) GetOrderItemsByOrderID(ctx context.Context, tx Transaction, orderID int64) ([]OrderItem, error) {
	orderItemList := make([]OrderItem, 0)

	// queryExecutor := ods.initiateQueryExecutor(tx)
	// err := queryExecutor.Where("order_id = ?", orderID).Find(&orderItemList).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return orderItemList, err
	// }

	return orderItemList, nil
}

func (ods *orderItemStore) StoreOrderItems(ctx context.Context, tx Transaction, orderItems []OrderItem) error {
	// queryExecutor := ods.initiateQueryExecutor(tx)
	// err := queryExecutor.Create(&orderItems).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}

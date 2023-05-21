package repository

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type orderItemStore struct {
	BaseRepository
}

func NewOrderItemRepo(db *storm.DB) repository.OrderItemStorer {
	return &orderItemStore{
		BaseRepository: BaseRepository{db},
	}
}

func (ods *orderItemStore) GetOrderItemsByOrderID(ctx context.Context, tx repository.Transaction, orderID int64) ([]repository.OrderItem, error) {
	orderItemList := make([]repository.OrderItem, 0)

	queryExecutor := ods.initiateQueryExecutor(tx)
	err := queryExecutor.Find("OrderID", orderID, &orderItemList)
	if err != nil && err != storm.ErrNotFound {
		return orderItemList, err
	}

	return orderItemList, nil
}

func (ods *orderItemStore) StoreOrderItems(ctx context.Context, tx repository.Transaction, orderItems []repository.OrderItem) error {
	queryExecutor := ods.initiateQueryExecutor(tx)
	for _, orderItem := range orderItems {

		//setting time fields
		orderItem.CreatedAt = ods.TimeNow()
		orderItem.UpdatedAt = ods.TimeNow()

		err := queryExecutor.Save(&orderItem)
		if err != nil {
			return err
		}
	}

	return nil
}

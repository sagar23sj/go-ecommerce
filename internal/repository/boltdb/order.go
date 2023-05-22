package repository

import (
	"context"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type orderStore struct {
	BaseRepository
}

func NewOrderRepo(db *storm.DB) repository.OrderStorer {
	return &orderStore{
		BaseRepository: BaseRepository{db},
	}
}

func (os *orderStore) GetOrderByID(ctx context.Context, tx repository.Transaction, orderID int64) (repository.Order, error) {
	var order repository.Order

	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.One("ID", orderID, &order)
	if err != nil && err != storm.ErrNotFound {
		return repository.Order{}, err
	}

	return order, nil
}

func (os *orderStore) CreateOrder(ctx context.Context, tx repository.Transaction, order repository.Order) (repository.Order, error) {

	queryExecutor := os.initiateQueryExecutor(tx)

	order.CreatedAt = os.TimeNow()
	order.UpdatedAt = os.TimeNow()
	err := queryExecutor.Save(&order)
	if err != nil {
		return repository.Order{}, err
	}

	return order, nil
}

func (os *orderStore) UpdateOrderStatus(ctx context.Context, tx repository.Transaction, orderID int64, status string) error {
	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.Update(&repository.Order{ID: uint(orderID), Status: status, UpdatedAt: os.TimeNow()})
	if err != nil {
		return err
	}

	return nil
}

func (os *orderStore) UpdateOrderDispatchDate(ctx context.Context, tx repository.Transaction, orderID int64, dispatchedAt time.Time) error {
	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.Update(&repository.Order{ID: uint(orderID), DispatchedAt: dispatchedAt, UpdatedAt: os.TimeNow()})
	if err != nil {
		return err
	}

	return nil
}

func (os *orderStore) ListOrders(ctx context.Context, tx repository.Transaction) ([]repository.Order, error) {
	orderList := make([]repository.Order, 0)

	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.All(&orderList)
	if err != nil {
		return orderList, err
	}

	return orderList, nil
}

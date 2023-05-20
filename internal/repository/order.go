package repository

import (
	"context"
	"time"

	"github.com/asdine/storm/v3"
)

type orderStore struct {
	BaseRepository
}

type OrderStorer interface {
	RepositoryTransaction

	GetOrderByID(ctx context.Context, tx Transaction, orderID int64) (Order, error)
	CreateOrder(ctx context.Context, tx Transaction, order Order) (Order, error)
	UpdateOrderStatus(ctx context.Context, tx Transaction, orderID int64, status string) error
	UpdateOrderDispatchDate(ctx context.Context, tx Transaction, orderID int64, dispatchedAt time.Time) error
	ListOrders(ctx context.Context, tx Transaction) ([]Order, error)
}

func NewOrderRepo(db *storm.DB) OrderStorer {
	return &orderStore{
		BaseRepository: BaseRepository{db},
	}
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

func (os *orderStore) GetOrderByID(ctx context.Context, tx Transaction, orderID int64) (Order, error) {
	var order Order

	// queryExecutor := os.initiateQueryExecutor(tx)
	// err := queryExecutor.First(&order, orderID).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return Order{}, err
	// }

	return order, nil
}

func (os *orderStore) CreateOrder(ctx context.Context, tx Transaction, order Order) (Order, error) {
	// queryExecutor := os.initiateQueryExecutor(tx)
	// err := queryExecutor.Create(&order).Error
	// if err != nil {
	// 	return Order{}, err
	// }

	return order, nil
}

func (os *orderStore) UpdateOrderStatus(ctx context.Context, tx Transaction, orderID int64, status string) error {
	// queryExecutor := os.initiateQueryExecutor(tx)
	// err := queryExecutor.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (os *orderStore) UpdateOrderDispatchDate(ctx context.Context, tx Transaction, orderID int64, dispatchedAt time.Time) error {
	// queryExecutor := os.initiateQueryExecutor(tx)
	// err := queryExecutor.Model(&Order{}).Where("id = ?", orderID).Update("dispatched_at", dispatchedAt).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (os *orderStore) ListOrders(ctx context.Context, tx Transaction) ([]Order, error) {
	orderList := make([]Order, 0)

	// queryExecutor := os.initiateQueryExecutor(tx)
	// err := queryExecutor.Find(&orderList).Error
	// if err != nil {
	// 	return orderList, err
	// }

	return orderList, nil
}

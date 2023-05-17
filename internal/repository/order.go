package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type orderStore struct {
	BaseRepository
}

type OrderStorer interface {
	RepositoryTransaction

	GetOrderByID(ctx context.Context, tx *gorm.DB, orderID int64) (Order, error)
	CreateOrder(ctx context.Context, tx *gorm.DB, order Order) (Order, error)
	UpdateOrderStatus(ctx context.Context, tx *gorm.DB, orderID int64, status string) error
	UpdateOrderDispatchDate(ctx context.Context, tx *gorm.DB, orderID int64, dispatchedAt time.Time) error
	ListOrders(ctx context.Context, tx *gorm.DB) ([]Order, error)
}

func NewOrderRepo(db *gorm.DB) OrderStorer {
	return &orderStore{
		BaseRepository: BaseRepository{db},
	}
}

type Order struct {
	gorm.Model
	Amount             float64
	DiscountPercentage float64
	FinalAmount        float64
	Status             string
	DispatchedAt       time.Time
}

func (os *orderStore) GetOrderByID(ctx context.Context, tx *gorm.DB, orderID int64) (Order, error) {
	var order Order

	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.First(&order, orderID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return Order{}, err
	}

	return order, nil
}

func (os *orderStore) CreateOrder(ctx context.Context, tx *gorm.DB, order Order) (Order, error) {
	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.Create(&order).Error
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (os *orderStore) UpdateOrderStatus(ctx context.Context, tx *gorm.DB, orderID int64, status string) error {
	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}

func (os *orderStore) UpdateOrderDispatchDate(ctx context.Context, tx *gorm.DB, orderID int64, dispatchedAt time.Time) error {
	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.Model(&Order{}).Where("id = ?", orderID).Update("dispatched_at", dispatchedAt).Error
	if err != nil {
		return err
	}

	return nil
}

func (os *orderStore) ListOrders(ctx context.Context, tx *gorm.DB) ([]Order, error) {
	orderList := make([]Order, 0)

	queryExecutor := os.initiateQueryExecutor(tx)
	err := queryExecutor.Find(&orderList).Error
	if err != nil {
		return orderList, err
	}

	return orderList, nil
}

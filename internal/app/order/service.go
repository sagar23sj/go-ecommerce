package order

import (
	"context"

	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type service struct {
	orderRepo repository.OrderStorer
}

type Service interface {
	CreateOrder(ctx context.Context, orderDetails Order) (Order, error)
	GetOrderDetailsByID(ctx context.Context, orderID int64) (Order, error)
	ListOrders(ctx context.Context) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) (Order, error)
}

func NewService(orderRepo repository.OrderStorer) Service {
	return &service{
		orderRepo: orderRepo,
	}
}

func (os *service) CreateOrder(ctx context.Context, orderDetails Order) (Order, error) {
	return Order{}, nil
}

func (os *service) GetOrderDetailsByID(ctx context.Context, orderID int64) (Order, error) {
	return Order{}, nil
}

func (os *service) ListOrders(ctx context.Context) ([]Order, error) {
	orderList := make([]Order, 0)
	return orderList, nil
}

func (os *service) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (Order, error) {
	return Order{}, nil
}

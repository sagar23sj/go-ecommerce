package order

import (
	"context"

	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type service struct {
	orderRepo      repository.OrderStorer
	orderItemsRepo repository.OrderItemStorer
	productRepo    repository.ProductStorer
}

type Service interface {
	CreateOrder(ctx context.Context, orderDetails dto.CreateOrderRequest) (dto.Order, error)
	GetOrderDetailsByID(ctx context.Context, orderID int64) (dto.Order, error)
	ListOrders(ctx context.Context) ([]dto.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) (dto.Order, error)
}

func NewService(orderRepo repository.OrderStorer, orderItemsRepo repository.OrderItemStorer,
	productRepo repository.ProductStorer) Service {
	return &service{
		orderRepo:      orderRepo,
		orderItemsRepo: orderItemsRepo,
		productRepo:    productRepo,
	}
}

func (os *service) CreateOrder(ctx context.Context, orderDetails dto.CreateOrderRequest) (dto.Order, error) {
	return dto.Order{}, nil
}

func (os *service) GetOrderDetailsByID(ctx context.Context, orderID int64) (dto.Order, error) {
	return dto.Order{}, nil
}

func (os *service) ListOrders(ctx context.Context) ([]dto.Order, error) {
	orderList := make([]dto.Order, 0)
	return orderList, nil
}

func (os *service) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (dto.Order, error) {
	return dto.Order{}, nil
}

package order

import (
	"context"
	"fmt"
	"time"

	"github.com/sagar23sj/go-ecommerce/internal/app/product"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

var now = time.Now

type service struct {
	orderRepo      repository.OrderStorer
	orderItemsRepo repository.OrderItemStorer
	productSvc     product.Service
}

type Service interface {
	CreateOrder(ctx context.Context, orderDetails dto.CreateOrderRequest) (dto.Order, error)
	GetOrderDetailsByID(ctx context.Context, orderID int64) (dto.Order, error)
	ListOrders(ctx context.Context) ([]dto.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) (dto.Order, error)
}

func NewService(orderRepo repository.OrderStorer, orderItemsRepo repository.OrderItemStorer,
	productSvc product.Service) Service {
	return &service{
		orderRepo:      orderRepo,
		orderItemsRepo: orderItemsRepo,
		productSvc:     productSvc,
	}
}

func (os *service) CreateOrder(ctx context.Context, orderDetails dto.CreateOrderRequest) (order dto.Order, err error) {
	//initializing database transaction
	tx, err := os.orderRepo.BeginTx(ctx)
	if err != nil {
		return dto.Order{}, err
	}

	defer func() {
		txErr := os.orderRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	orderRepoObj, updatedProductInfo, err := os.calculateOrderValueFromProducts(ctx, tx, orderDetails.Products)
	if err != nil {
		return dto.Order{}, err
	}

	//Set Order Status to Placed
	orderRepoObj.Status = ListOrderStatus[OrderPlaced]

	//1. Inserting Order in Database
	orderDB, err := os.orderRepo.CreateOrder(ctx, tx, orderRepoObj)
	if err != nil {
		return dto.Order{}, err
	}

	orderItems := make([]repository.OrderItem, 0)
	for _, p := range orderDetails.Products {
		orderItems = append(orderItems, repository.OrderItem{
			OrderID:   int64(orderDB.ID),
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
		})
	}

	//2. Inserting order items in database
	err = os.orderItemsRepo.StoreOrderItems(ctx, tx, orderItems)
	if err != nil {
		return dto.Order{}, err
	}

	//3. Update Product quantity in database
	productQuantityMap := make(map[int64]int64)
	for _, p := range updatedProductInfo {
		productQuantityMap[p.ProductID] = p.Quantity
	}

	err = os.productSvc.UpdateProductQuantity(ctx, tx, productQuantityMap)
	if err != nil {
		return dto.Order{}, err
	}

	order = MapOrderRepoToOrderDto(orderDB, orderItems...)
	return order, nil
}

func (os *service) GetOrderDetailsByID(ctx context.Context, orderID int64) (order dto.Order, err error) {
	orderInfoDB, err := os.orderRepo.GetOrderByID(ctx, nil, orderID)
	if err != nil {
		return dto.Order{}, err
	}

	if orderInfoDB.ID == 0 {
		return dto.Order{}, apperrors.OrderNotFound{ID: orderID}
	}

	orderItemsDB, err := os.orderItemsRepo.GetOrderItemsByOrderID(ctx, nil, orderID)
	if err != nil {
		return dto.Order{}, err
	}

	order = MapOrderRepoToOrderDto(orderInfoDB, orderItemsDB...)
	return order, nil
}

func (os *service) ListOrders(ctx context.Context) ([]dto.Order, error) {
	orderList := make([]dto.Order, 0)

	orderListDB, err := os.orderRepo.ListOrders(ctx, nil)
	if err != nil {
		return orderList, err
	}

	for _, order := range orderListDB {
		orderList = append(orderList, MapOrderRepoToOrderDto(order))
	}

	return orderList, nil
}

func (os *service) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (order dto.Order, err error) {
	//initializing database transaction
	tx, err := os.orderRepo.BeginTx(ctx)
	if err != nil {
		return dto.Order{}, err
	}

	defer func() {
		txErr := os.orderRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	//order status invalid, return error OrderStatusInvalid
	if _, ok := MapOrderStatus[status]; !ok {
		return dto.Order{}, apperrors.OrderStatusInvalid{ID: orderID}
	}

	orderInfoDB, err := os.orderRepo.GetOrderByID(ctx, tx, orderID)
	if err != nil {
		return dto.Order{}, err
	}

	//order not found invalid, return error OrderNotFound
	if orderInfoDB.ID == 0 {
		return dto.Order{}, apperrors.OrderNotFound{ID: orderID}
	}

	//order status not allowed for update, return error OrderUpdationInvalid
	isUpdationValid := validateUpdateOrderStatusRequest(status, orderInfoDB.Status)
	if !isUpdationValid {
		return dto.Order{}, apperrors.OrderUpdationInvalid{
			ID:             orderID,
			CurrentState:   orderInfoDB.Status,
			RequestedState: status,
		}
	}

	//update order status in db
	err = os.orderRepo.UpdateOrderStatus(ctx, tx, orderID, status)
	if err != nil {
		return dto.Order{}, fmt.Errorf("error occured while updating order status: %w", err)
	}

	//update product quantity if order cancelled or returned
	if (MapOrderStatus[status] == OrderCancelled) || (MapOrderStatus[status] == OrderReturned) {

		orderItemsDB, err := os.orderItemsRepo.GetOrderItemsByOrderID(ctx, tx, orderID)
		if err != nil {
			return dto.Order{}, fmt.Errorf("error occured while fetching order items: %w", err)
		}

		productQuantityMap := make(map[int64]int64)
		for _, item := range orderItemsDB {

			product, err := os.productSvc.GetProductByID(ctx, tx, item.ProductID)
			if err != nil {
				return dto.Order{}, fmt.Errorf("error occured while fetching product with id %d,  %w", item.ProductID, err)
			}

			productQuantityMap[item.ProductID] = product.Quantity + item.Quantity
		}

		err = os.productSvc.UpdateProductQuantity(ctx, tx, productQuantityMap)
		if err != nil {
			return dto.Order{}, fmt.Errorf("error occured while updating product quantiry,  %w", err)
		}

	}

	//update dispatch date only when order_status = Dispatched
	if MapOrderStatus[status] == OrderDispatched {
		orderDispatchedAt := now()
		err = os.orderRepo.UpdateOrderDispatchDate(ctx, tx, orderID, orderDispatchedAt)
		if err != nil {
			return dto.Order{}, fmt.Errorf("error occured while updating order dispatch date: %w", err)
		}
	}

	orderInfoDB, err = os.orderRepo.GetOrderByID(ctx, tx, orderID)
	if err != nil {
		return dto.Order{}, err
	}

	order = MapOrderRepoToOrderDto(orderInfoDB)
	return order, err
}

func (os *service) calculateOrderValueFromProducts(ctx context.Context, tx repository.Transaction, requestedProducts []dto.ProductInfo) (
	orderInfo repository.Order, productsUpdated []dto.ProductInfo, err error) {

	productsUpdated = make([]dto.ProductInfo, 0)
	premiumProductCount := 0

	var orderAmount float64
	var discountPercent float64
	var finalOrderAmount float64

	for _, p := range requestedProducts {
		productInfo, err := os.productSvc.GetProductByID(ctx, tx, p.ProductID)
		if err != nil {
			return repository.Order{}, productsUpdated, err
		}

		//product quantity exceeded limit, return error apperrors.ProductQuantityExceeded
		if p.Quantity > MaxProductQuantity {
			return repository.Order{}, productsUpdated, apperrors.ProductQuantityExceeded{
				ID:            p.ProductID,
				QuantityAsked: p.Quantity,
				QuantityLimit: MaxProductQuantity,
			}
		}

		//product quantity insufficient, return error apperrors.ProductQuantityInsufficient
		if productInfo.Quantity < p.Quantity {
			return repository.Order{}, productsUpdated, apperrors.ProductQuantityInsufficient{
				ID:                p.ProductID,
				QuantityAsked:     p.Quantity,
				QuantityRemaining: productInfo.Quantity,
			}
		}

		orderAmount = orderAmount + (float64(p.Quantity) * productInfo.Price)

		//update premium product counter
		if productInfo.Category == string(product.PremiumProduct) {
			premiumProductCount = premiumProductCount + 1
		}

		//adding product details with updated quantity to the list
		productsUpdated = append(productsUpdated, dto.ProductInfo{
			ProductID: p.ProductID,
			Quantity:  (productInfo.Quantity - p.Quantity),
		})
	}

	finalOrderAmount = orderAmount
	//checking if premium products are equal or more than 3
	if premiumProductCount >= PremiumProductsForDiscount {
		discountPercent = DefaultDiscountPercentage
		finalOrderAmount = orderAmount * (100 - discountPercent) / 100
	}

	orderInfo = repository.Order{
		Amount:             orderAmount,
		DiscountPercentage: discountPercent,
		FinalAmount:        finalOrderAmount,
	}

	return orderInfo, productsUpdated, nil
}

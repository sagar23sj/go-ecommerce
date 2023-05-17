package order

import (
	"context"
	"fmt"
	"time"

	"github.com/sagar23sj/go-ecommerce/internal/app/product"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
	"gorm.io/gorm"
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

	//Set Order Status to Default Placed
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

	err = os.productRepo.UpdateProductQuantity(ctx, tx, productQuantityMap)
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
	isUpdationValid := os.validateUpdateOrderStatusRequest(ctx, status, orderInfoDB.Status)
	if !isUpdationValid {
		return dto.Order{}, apperrors.OrderUpdationInvalid{
			ID:             orderID,
			CurrentState:   orderInfoDB.Status,
			RequestedState: status,
		}
	}

	err = os.orderRepo.UpdateOrderStatus(ctx, tx, orderID, status)
	if err != nil {
		return dto.Order{}, fmt.Errorf("error occured while updating order status: %w", err)
	}

	//update dispatch date only when order is dispatched
	if MapOrderStatus[status] == OrderDispatched {
		orderDispatchedAt := time.Now()
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

func (os *service) calculateOrderValueFromProducts(ctx context.Context, tx *gorm.DB, requestedProducts []dto.ProductInfo) (
	orderInfo repository.Order, productsUpdated []dto.ProductInfo, err error) {

	productsUpdated = make([]dto.ProductInfo, 0)
	premiumProductCount := 0

	var orderAmount float64
	var discountPercent float64
	var finalOrderAmount float64

	//merging multiple product with same ID
	productQuantityMap := make(map[int64]int64)
	for _, p := range requestedProducts {
		if _, ok := productQuantityMap[p.ProductID]; !ok {
			productQuantityMap[p.ProductID] = p.Quantity
		}

		productQuantityMap[p.ProductID] = productQuantityMap[p.ProductID] + p.Quantity
	}

	for productID, productQuantity := range productQuantityMap {
		productInfoDB, err := os.productRepo.GetProductByID(ctx, tx, productID)
		if err != nil {
			return repository.Order{}, productsUpdated, err
		}

		//product not found, return error apperrors.ProductNotFound
		if productInfoDB.ID == 0 {
			return repository.Order{}, productsUpdated, apperrors.ProductNotFound{ID: productID}
		}

		//product quantity insufficient, return error apperrors.ProductQuantityInsufficient
		if productInfoDB.Quantity < productQuantity {
			return repository.Order{}, productsUpdated, apperrors.ProductQuantityInsufficient{
				ID:                productID,
				QuantityAsked:     productQuantity,
				QuantityRemaining: productInfoDB.Quantity,
			}
		}

		//product quantity exceeded limit, return error apperrors.ProductQuantityExceeded
		if productQuantity > product.MaxProductQuantity {
			return repository.Order{}, productsUpdated, apperrors.ProductQuantityExceeded{
				ID:            productID,
				QuantityAsked: productQuantity,
				QuantityLimit: product.MaxProductQuantity,
			}
		}

		orderAmount = orderAmount + (float64(productQuantity) * productInfoDB.Price)

		//update premium product counter
		if productInfoDB.Category == string(product.PremiumProduct) {
			premiumProductCount = premiumProductCount + 1
		}

		//adding product details with updated quantity to the list
		productsUpdated = append(productsUpdated, dto.ProductInfo{
			ProductID: productID,
			Quantity:  (productInfoDB.Quantity - productQuantity),
		})
	}

	finalOrderAmount = orderAmount
	//checking if premium products are equal or more than 3
	if premiumProductCount >= product.PremiumProductsForDiscount {
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

func (os *service) validateUpdateOrderStatusRequest(ctx context.Context, RequestOrderStatus, DBOrderStatus string) (isUpdateValid bool) {
	requestedOrderState := MapOrderStatus[RequestOrderStatus]
	currentOrderState := MapOrderStatus[DBOrderStatus]

	//donot update if requested and current state is same
	if currentOrderState == requestedOrderState {
		return false
	}

	//donot update if order is already cancelled
	if currentOrderState == OrderCancelled {
		return false
	}

	//allow cancel only before order is completed
	if requestedOrderState == OrderCancelled && currentOrderState < OrderCompleted {
		return true
	}

	//order state update should not go backwards unless it is cancel reqeust
	if requestedOrderState < currentOrderState {
		return false
	}

	//order status update can only go one step forward
	if requestedOrderState != (currentOrderState + 1) {
		return false
	}

	return true
}

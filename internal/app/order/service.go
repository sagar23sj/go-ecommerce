package order

import (
	"context"

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

	order = dto.Order{
		ID:                 int64(orderDB.ID),
		Products:           orderDetails.Products,
		Amount:             orderDB.Amount,
		DiscountPercentage: orderDB.DiscountPercentage,
		DiscountedAmount:   orderDB.DiscountedAmount,
		Status:             orderDB.Status,
		CreatedAt:          orderDB.CreatedAt,
		UpdatedAt:          orderDB.UpdatedAt,
	}

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

	products := make([]dto.ProductInfo, 0)
	for _, items := range orderItemsDB {
		products = append(products, dto.ProductInfo{
			ProductID: items.ProductID,
			Quantity:  items.Quantity,
		})
	}

	order = dto.Order{
		ID:                 int64(orderInfoDB.ID),
		Products:           products,
		Amount:             orderInfoDB.Amount,
		DiscountPercentage: orderInfoDB.DiscountPercentage,
		DiscountedAmount:   orderInfoDB.DiscountedAmount,
		Status:             orderInfoDB.Status,
		CreatedAt:          orderInfoDB.CreatedAt,
		UpdatedAt:          orderInfoDB.UpdatedAt,
	}

	return order, nil
}

func (os *service) ListOrders(ctx context.Context) ([]dto.Order, error) {
	orderList := make([]dto.Order, 0)

	orderListDB, err := os.orderRepo.ListOrders(ctx, nil)
	if err != nil {
		return orderList, err
	}

	for _, order := range orderListDB {
		orderList = append(orderList, dto.Order{
			ID:                 int64(order.ID),
			Amount:             order.Amount,
			DiscountPercentage: order.DiscountPercentage,
			DiscountedAmount:   order.DiscountedAmount,
			Status:             order.Status,
			CreatedAt:          order.CreatedAt,
			UpdatedAt:          order.UpdatedAt,
		})
	}

	return orderList, nil
}

func (os *service) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (order dto.Order, err error) {
	if _, ok := MapOrderStatus[status]; !ok {
		return dto.Order{}, apperrors.OrderStatusInvalid{ID: orderID}
	}

	orderInfo, err := os.orderRepo.GetOrderByID(ctx, nil, orderID)
	if err != nil {
		return dto.Order{}, err
	}

	if orderInfo.ID == 0 {
		return dto.Order{}, apperrors.OrderNotFound{ID: orderID}
	}

	if MapOrderStatus[orderInfo.Status] <= MapOrderStatus[status] {
		return dto.Order{}, apperrors.OrderStatusInvalid{ID: orderID}
	}

	return dto.Order{}, nil
}

func (os *service) calculateOrderValueFromProducts(ctx context.Context, tx *gorm.DB, requestedProducts []dto.ProductInfo) (
	orderInfo repository.Order, productsUpdated []dto.ProductInfo, err error) {

	productsUpdated = make([]dto.ProductInfo, 0)
	premiumProductCount := 0

	var orderAmount float64
	var discountPercent float64
	var discountedOrderAmount float64

	for _, p := range requestedProducts {
		productInfo, err := os.productRepo.GetProductByID(ctx, tx, p.ProductID)
		if err != nil {
			return repository.Order{}, productsUpdated, err
		}

		if productInfo.ID == 0 {
			return repository.Order{}, productsUpdated, apperrors.ProductNotFound{ID: int64(p.ProductID)}
		}

		if productInfo.Quantity < p.Quantity {
			return repository.Order{}, productsUpdated, apperrors.ProductQuantityInsufficient{
				ID:                p.ProductID,
				QuantityAsked:     p.Quantity,
				QuantityRemaining: productInfo.Quantity,
			}
		}

		if p.Quantity > product.MaxProductQuantity {
			return repository.Order{}, productsUpdated, apperrors.ProductQuantityExceeded{
				ID:            p.ProductID,
				QuantityAsked: p.Quantity,
				QuantityLimit: product.MaxProductQuantity,
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

	//checking if premium products are equal or more than 3
	if premiumProductCount >= product.PremiumProductsForDiscount {
		discountPercent = DiscountPercentage
		discountedOrderAmount = orderAmount * (100 - discountPercent) / 100
	}

	orderInfo = repository.Order{
		Amount:             orderAmount,
		DiscountPercentage: discountPercent,
		DiscountedAmount:   discountedOrderAmount,
	}

	return orderInfo, productsUpdated, nil
}

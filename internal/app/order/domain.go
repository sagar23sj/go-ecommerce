package order

import (
	"time"

	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

const (
	DefaultDiscountPercentage  = 10
	MaxProductQuantity         = 10
	PremiumProductsForDiscount = 3
)

type OrderStatus int

const (
	OrderCancelled OrderStatus = iota
	OrderPlaced
	OrderDispatched
	OrderCompleted
	OrderReturned
)

var MapOrderStatus = map[string]OrderStatus{
	"Cancelled":  OrderCancelled,
	"Placed":     OrderPlaced,
	"Dispatched": OrderDispatched,
	"Completed":  OrderCompleted,
	"Returned":   OrderReturned,
}

// Note -- the order of this slice needs to match
// the order of the iota enum values defined above
var ListOrderStatus = []string{
	"Cancelled",
	"Placed",
	"Dispatched",
	"Completed",
	"Returned",
}

func validateUpdateOrderStatusRequest(RequestOrderStatus, DBOrderStatus string) (isUpdateValid bool) {
	requestedOrderState := MapOrderStatus[RequestOrderStatus]
	currentOrderState := MapOrderStatus[DBOrderStatus]

	//donot update if order is already cancelled
	if currentOrderState == OrderCancelled {
		return false
	}

	//allow cancel only before order is completed
	if requestedOrderState == OrderCancelled && currentOrderState < OrderCompleted {
		return true
	}

	//donot update if requested state is same or lower to current state
	if currentOrderState >= requestedOrderState {
		return false
	}

	//order status update can only go one step forward
	if requestedOrderState != (currentOrderState + 1) {
		return false
	}

	return true
}

func MapOrderRepoToOrderDto(order repository.Order, orderItems ...repository.OrderItem) dto.Order {

	productInfo := make([]dto.ProductInfo, 0)
	for _, orderItem := range orderItems {
		productInfo = append(productInfo, dto.ProductInfo{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		})
	}

	var dispatchedAt *time.Time = &order.DispatchedAt
	if order.DispatchedAt.IsZero() {
		dispatchedAt = nil
	}

	return dto.Order{
		ID:                 int64(order.ID),
		Products:           productInfo,
		Amount:             order.Amount,
		DiscountPercentage: order.DiscountPercentage,
		FinalAmount:        order.FinalAmount,
		Status:             order.Status,
		DispatchedAt:       dispatchedAt,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}
}

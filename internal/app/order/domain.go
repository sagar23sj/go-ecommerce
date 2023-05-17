package order

import (
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

const (
	DefaultDiscountPercentage = 10
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

func MapOrderRepoToOrderDto(order repository.Order, orderItems ...repository.OrderItem) dto.Order {

	productInfo := make([]dto.ProductInfo, 0)
	for _, orderItem := range orderItems {
		productInfo = append(productInfo, dto.ProductInfo{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		})
	}

	return dto.Order{
		ID:                 int64(order.ID),
		Products:           productInfo,
		Amount:             order.Amount,
		DiscountPercentage: order.DiscountPercentage,
		FinalAmount:        order.FinalAmount,
		Status:             order.Status,
		DispatchedAt:       order.DispatchedAt,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}
}

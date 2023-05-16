package order

import (
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

func MapOrderRepoToOrderDto(order repository.Order, orderItems []repository.OrderItem) dto.Order {

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
		DiscountedAmount:   order.DiscountedAmount,
		Status:             order.Status,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}
}

package app

import (
	"github.com/sagar23sj/go-ecommerce/internal/app/order"
	"github.com/sagar23sj/go-ecommerce/internal/app/product"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
	"gorm.io/gorm"
)

type Dependencies struct {
	OrderService   order.Service
	ProductService product.Service
}

func NewServices(db *gorm.DB) Dependencies {
	//initialize repo dependencies
	orderRepo := repository.NewOrderRepo(db)
	orderItemsRepo := repository.NewOrderItemRepo(db)
	productRepo := repository.NewProductRepo(db)

	//initialize service dependencies
	productService := product.NewService(productRepo)
	orderService := order.NewService(orderRepo, orderItemsRepo, productService)

	return Dependencies{
		OrderService:   orderService,
		ProductService: productService,
	}
}

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
	productRepo := repository.NewProductRepo(db)

	//initialize service dependencies
	orderService := order.NewService(orderRepo)
	productService := product.NewService(productRepo)

	return Dependencies{
		OrderService:   orderService,
		ProductService: productService,
	}
}

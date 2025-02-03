package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sagar23sj/go-ecommerce/internal/app/order"
	"github.com/sagar23sj/go-ecommerce/internal/app/product"
)

func NewRouter(deps Dependencies) chi.Router {
	router := chi.NewRouter()

	//order APIs
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Post("/order", order.CreateOrderHandler(deps.OrderService))
		r.Get("/orders", order.ListOrdersHandler(deps.OrderService))
		r.Get("/order/{id}", order.GetOrderDetailsHandler(deps.OrderService))
		r.Patch("/order/{id}/status", order.UpdateOrderStatusHandler(deps.OrderService))

	})

	//product APIs
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/product/{id}", product.GetProductHandler(deps.ProductService))
		r.Get("/products", product.ListProductHandler(deps.ProductService))

	})

	return router
}

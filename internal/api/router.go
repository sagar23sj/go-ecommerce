package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	orderApi "github.com/sagar23sj/go-ecommerce/internal/api/order"
	productApi "github.com/sagar23sj/go-ecommerce/internal/api/product"
	"github.com/sagar23sj/go-ecommerce/internal/app"
)

func NewRouter(deps app.Dependencies) chi.Router {
	router := chi.NewRouter()

	//order APIs
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Post("/order", orderApi.CreateOrderHandler(deps.OrderService))
		r.Get("/order/{id}", orderApi.GetOrderHandler(deps.OrderService))
		r.Patch("/order/{id}/status", orderApi.UpdateOrderStatusHandler(deps.OrderService))

	})

	//product APIs
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/product/{id}", productApi.GetProductHandler(deps.ProductService))
		r.Get("/products", productApi.ListProductHandler(deps.ProductService))

	})

	return router
}

package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sagar23sj/go-ecommerce/internal/app"
)

func NewRouter(deps app.Dependencies) chi.Router {
	router := chi.NewRouter()

	//order APIs
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Post("/order", CreateOrderHandler(deps.OrderService))
		r.Get("/order/{id}", GetOrderHandler(deps.OrderService))
		r.Patch("/order/{id}/status", UpdateOrderStatusHandler(deps.OrderService))

	})

	//product APIs
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/product/{id}", getProductHandler(deps.ProductService))
		r.Get("/products", listProductHandler(deps.ProductService))

	})

	return router
}

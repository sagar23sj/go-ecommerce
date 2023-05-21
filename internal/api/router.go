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

		r.Post("/orders", createOrderHandler(deps.OrderService))
		r.Get("/orders", listOrdersHandler(deps.OrderService))
		r.Get("/orders/{id}", getOrderDetailsHandler(deps.OrderService))
		r.Patch("/orders/{id}/status", updateOrderStatusHandler(deps.OrderService))

	})

	//product APIs
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/products/{id}", getProductHandler(deps.ProductService))
		r.Get("/products", listProductHandler(deps.ProductService))

	})

	return router
}

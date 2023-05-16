package api

import (
	"encoding/json"
	"net/http"

	"github.com/sagar23sj/go-ecommerce/internal/app/order"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/logger"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/middleware"
	"go.uber.org/zap"
)

func createOrderHandler(orderSvc order.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req dto.CreateOrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Errorw(ctx, "error occured while decoding request",
				zap.Error(err),
			)
			middleware.ErrorResponse(ctx, w, http.StatusBadRequest, apperrors.ErrInvalidRequestBody)
			return
		}

		orderInfo, err := orderSvc.CreateOrder(ctx, req)
		if err != nil {
			logger.Errorw(ctx, "error occured while creating order",
				zap.Error(err),
			)
			middleware.ErrorResponse(ctx, w, http.StatusInternalServerError, apperrors.ErrInternalServerError)
			return
		}

		middleware.SuccessResponse(ctx, w, http.StatusCreated, orderInfo)
	}
}

func getOrderHandler(orderSvc order.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func updateOrderStatusHandler(orderSvc order.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

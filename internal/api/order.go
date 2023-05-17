package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
			statusCode, err := apperrors.MapError(err)
			middleware.ErrorResponse(ctx, w, statusCode, err)
			return
		}

		middleware.SuccessResponse(ctx, w, http.StatusCreated, orderInfo)
	}
}

func getOrderHandler(orderSvc order.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rawOrderID := chi.URLParam(r, "id")
		orderID, err := strconv.Atoi(rawOrderID)
		if err != nil {
			logger.Errorw(ctx, "error occured while converting orderID to an integer",
				zap.Error(err),
				zap.String("id", rawOrderID),
			)

			middleware.ErrorResponse(ctx, w, http.StatusInternalServerError, apperrors.ErrInvalidRequestParam)
			return
		}

		response, err := orderSvc.GetOrderDetailsByID(ctx, int64(orderID))
		if err != nil {
			logger.Errorw(ctx, "error occured while fetching order info",
				zap.Error(err),
			)

			statusCode, errResponse := apperrors.MapError(err)
			middleware.ErrorResponse(ctx, w, statusCode, errResponse)
			return
		}

		middleware.SuccessResponse(ctx, w, http.StatusOK, response)
		return
	}
}

func updateOrderStatusHandler(orderSvc order.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req dto.UpdateOrderStatusRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Errorw(ctx, "error occured while decoding request",
				zap.Error(err),
			)
			middleware.ErrorResponse(ctx, w, http.StatusBadRequest, apperrors.ErrInvalidRequestBody)
			return
		}

		orderInfo, err := orderSvc.UpdateOrderStatus(ctx, req.OrderID, req.Status)
		if err != nil {
			logger.Errorw(ctx, "error occured while updating order status",
				zap.Error(err),
			)
			statusCode, err := apperrors.MapError(err)
			middleware.ErrorResponse(ctx, w, statusCode, err)
			return
		}

		middleware.SuccessResponse(ctx, w, http.StatusCreated, orderInfo)
	}
}

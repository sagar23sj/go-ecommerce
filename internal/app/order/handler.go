package order

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sagar23sj/go-ecommerce/internal/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/model"
	"github.com/sagar23sj/go-ecommerce/internal/utils"
)

func CreateOrderHandler(orderSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req model.CreateOrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while decoding request",
				slog.String("error", err.Error()),
			)
			utils.ErrorResponse(ctx, w, http.StatusBadRequest, apperrors.ErrInvalidRequestBody)
			return
		}

		err = req.Validate()
		if err != nil {
			slog.ErrorContext(ctx, "error occured while validating create order request",
				slog.String("error", err.Error()),
			)
			utils.ErrorResponse(ctx, w, http.StatusUnprocessableEntity, err)
			return
		}

		orderInfo, err := orderSvc.CreateOrder(ctx, req)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while creating order",
				slog.String("error", err.Error()),
			)
			statusCode, err := apperrors.MapError(err)
			utils.ErrorResponse(ctx, w, statusCode, err)
			return
		}

		utils.SuccessResponse(ctx, w, http.StatusCreated, orderInfo)
	}
}

func GetOrderDetailsHandler(orderSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rawOrderID := chi.URLParam(r, "id")
		orderID, err := strconv.Atoi(rawOrderID)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while converting orderID to an integer",
				slog.String("error", err.Error()),
				slog.String("id", rawOrderID),
			)

			utils.ErrorResponse(ctx, w, http.StatusBadRequest, apperrors.ErrInvalidRequestParam)
			return
		}

		response, err := orderSvc.GetOrderDetailsByID(ctx, int64(orderID))
		if err != nil {
			slog.ErrorContext(ctx, "error occured while fetching order info",
				slog.String("error", err.Error()),
			)

			statusCode, errResponse := apperrors.MapError(err)
			utils.ErrorResponse(ctx, w, statusCode, errResponse)
			return
		}

		utils.SuccessResponse(ctx, w, http.StatusOK, response)
		return
	}
}

func UpdateOrderStatusHandler(orderSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req model.UpdateOrderStatusRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while decoding request",
				slog.String("error", err.Error()),
			)
			utils.ErrorResponse(ctx, w, http.StatusBadRequest, apperrors.ErrInvalidRequestBody)
			return
		}

		orderInfo, err := orderSvc.UpdateOrderStatus(ctx, req.OrderID, req.Status)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while updating order status",
				slog.String("error", err.Error()),
			)
			statusCode, err := apperrors.MapError(err)
			utils.ErrorResponse(ctx, w, statusCode, err)
			return
		}

		utils.SuccessResponse(ctx, w, http.StatusCreated, orderInfo)
	}
}

func ListOrdersHandler(orderSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response, err := orderSvc.ListOrders(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while fetching orders list",
				slog.String("error", err.Error()),
			)

			utils.ErrorResponse(ctx, w, http.StatusInternalServerError, apperrors.ErrInternalServerError)
			return
		}

		utils.SuccessResponse(ctx, w, http.StatusOK, response)
		return
	}
}

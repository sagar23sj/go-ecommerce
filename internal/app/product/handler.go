package product

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sagar23sj/go-ecommerce/internal/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/utils"
)

func GetProductHandler(productSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rawProductID := chi.URLParam(r, "id")
		productID, err := strconv.Atoi(rawProductID)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while converting productID to an integer",
				slog.String("error", err.Error()),
				slog.String("id", rawProductID),
			)

			utils.ErrorResponse(ctx, w, http.StatusBadRequest, apperrors.ErrInvalidRequestParam)
			return
		}

		response, err := productSvc.GetProductByID(ctx, nil, int64(productID))
		if err != nil {
			slog.ErrorContext(ctx, "error occured while fetching product info",
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

func ListProductHandler(productSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response, err := productSvc.ListProducts(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "error occured while fetching product list",
				slog.String("error", err.Error()),
			)

			utils.ErrorResponse(ctx, w, http.StatusInternalServerError, apperrors.ErrInternalServerError)
			return
		}

		utils.SuccessResponse(ctx, w, http.StatusOK, response)
		return
	}
}

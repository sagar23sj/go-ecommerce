package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sagar23sj/go-ecommerce/internal/app/product"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/logger"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/middleware"
	"go.uber.org/zap"
)

func getProductHandler(productSvc product.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rawProductID := chi.URLParam(r, "id")
		productID, err := strconv.Atoi(rawProductID)
		if err != nil {
			logger.Errorw(ctx, "error occured while converting productID to an integer",
				zap.Error(err),
				zap.String("id", rawProductID),
			)

			middleware.ErrorResponse(ctx, w, http.StatusBadRequest, apperrors.ErrInvalidRequestParam)
			return
		}

		response, err := productSvc.GetProductByID(ctx, int64(productID))
		if err != nil {
			logger.Errorw(ctx, "error occured while fetching product info",
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

func listProductHandler(productSvc product.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response, err := productSvc.ListProducts(ctx)
		if err != nil {
			logger.Errorw(ctx, "error occured while fetching product list",
				zap.Error(err),
			)

			middleware.ErrorResponse(ctx, w, http.StatusInternalServerError, apperrors.ErrInternalServerError)
			return
		}

		middleware.SuccessResponse(ctx, w, http.StatusOK, response)
		return
	}
}

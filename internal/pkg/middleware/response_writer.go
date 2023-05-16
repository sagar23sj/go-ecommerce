package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sagar23sj/go-ecommerce/internal/pkg/logger"
	"go.uber.org/zap"
)

type response struct {
	ErrorCode    int         `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	Data         interface{} `json:"data"`
}

func SuccessResponse(ctx context.Context, w http.ResponseWriter, status int, data any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	payload := response{
		Data: data,
	}

	out, err := json.Marshal(payload)
	if err != nil {
		logger.Errorw(ctx, "cannot marshal success response payload", zap.Error(err))
		writeServerErrorResponse(ctx, w)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		logger.Errorw(ctx, "cannot write json success response", zap.Error(err))
		writeServerErrorResponse(ctx, w)
		return
	}
}

func ErrorResponse(ctx context.Context, w http.ResponseWriter, httpStatus int, err error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	payload := response{
		ErrorCode:    httpStatus,
		ErrorMessage: err.Error(),
	}

	out, err := json.Marshal(payload)
	if err != nil {
		logger.Errorw(ctx, "error occured while marshaling response payload", zap.Error(err))
		writeServerErrorResponse(ctx, w)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		logger.Errorw(ctx, "error occured while writing response", zap.Error(err))
		writeServerErrorResponse(ctx, w)
		return
	}
}

func writeServerErrorResponse(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(fmt.Sprintf("{\"message\":%s}", "internal server error")))
	if err != nil {
		logger.Errorw(ctx, "error occured while writing response", zap.Error(err))
	}
}

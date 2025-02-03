package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type response struct {
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
		slog.ErrorContext(ctx, "cannot marshal success response payload", slog.String("error", err.Error()))
		writeServerErrorResponse(ctx, w)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		slog.ErrorContext(ctx, "cannot write json success response", slog.String("error", err.Error()))
		writeServerErrorResponse(ctx, w)
		return
	}
}

func ErrorResponse(ctx context.Context, w http.ResponseWriter, httpStatus int, err error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	payload := response{
		ErrorMessage: err.Error(),
	}

	out, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "error occured while marshaling response payload", slog.String("error", err.Error()))
		writeServerErrorResponse(ctx, w)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		slog.ErrorContext(ctx, "error occured while writing response", slog.String("error", err.Error()))
		writeServerErrorResponse(ctx, w)
		return
	}
}

func writeServerErrorResponse(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(fmt.Sprintf("{\"message\":%s}", "internal server error")))
	if err != nil {
		slog.ErrorContext(ctx, "error occured while writing response", slog.String("error", err.Error()))
	}
}

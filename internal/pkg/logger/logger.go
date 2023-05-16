package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger *zap.SugaredLogger

func init() {
	zapLogger := getLogger()
	appLogger = zapLogger.Sugar()
}

func getLogger() (logger *zap.Logger) {

	logLevel := zapcore.InfoLevel
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, os.Stdout, logLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}

func Errorw(ctx context.Context, message string, args ...interface{}) {
	appLogger.Errorw(message, args...)
}

func Infow(ctx context.Context, message string, args ...interface{}) {
	appLogger.Infow(message, args...)
}

func Warnw(ctx context.Context, message string, args ...interface{}) {
	appLogger.Warnw(message, args...)
}

func Debugw(ctx context.Context, message string, args ...interface{}) {
	appLogger.Debugw(message, args...)
}

func Fatalw(ctx context.Context, message string, args ...interface{}) {
	appLogger.Fatalw(message, args...)
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/sagar23sj/go-ecommerce/internal/api"
	"github.com/sagar23sj/go-ecommerce/internal/app"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/constants"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/logger"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
	"go.uber.org/zap"
)

func main() {

	ctx := context.Background()
	logger.Infow(ctx, "Starting E-Commerce Application....")
	defer logger.Infow(ctx, "Shutting Down E-Commerce Application...")

	sqlDB, err := repository.InitializeDatabase()
	if err != nil {
		logger.Fatalw(ctx, "error occured while initializing database object",
			zap.Error(err),
		)
	}

	//initialize service dependencies
	services := app.NewServices(sqlDB)

	//initialize router
	router := api.NewRouter(services)

	var group run.Group
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", constants.HTTPPort),
		Handler: router,
	}

	//Adding HTTP Server to run group
	group.Add(
		func() error {
			{
				logger.Infow(ctx, "Starting HTTP Server", zap.Int("port", constants.HTTPPort))

				err := srv.ListenAndServe()
				if err != nil {
					logger.Errorw(ctx, "HTTP Server Closed", zap.Error(err))
				}
				return err
			}
		},
		func(err error) {
			logger.Infow(ctx, "Shutting HTTP server down gracefully...", zap.Error(err))

			ctx, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()

			err = srv.Shutdown(ctx)
			if err != nil {
				logger.Infow(ctx, "Cannot shut HTTP server down gracefully. Shutting it down forcefully...", zap.Error(err))
			}

			logger.Infow(ctx, "HTTP server shut down complete.")
		},
	)

	//Adding graceful shutdown handler to run group
	group.Add(
		run.SignalHandler(
			ctx,
			syscall.SIGABRT,
			syscall.SIGALRM,
			syscall.SIGBUS,
			syscall.SIGINT,
			syscall.SIGTERM,
		),
	)

	err = group.Run()
	if err != nil {
		logger.Infow(ctx, "Run group has been interrupted.", zap.Error(err))
	}

}

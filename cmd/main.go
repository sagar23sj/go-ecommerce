package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/sagar23sj/go-ecommerce/internal/app"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

const (
	HTTPPort = 8080
)

func main() {

	ctx := context.Background()
	slog.InfoContext(ctx, "Starting E-Commerce Application....")
	defer slog.InfoContext(ctx, "Shutting Down E-Commerce Application...")

	sqlDB, err := repository.InitializeDatabase()
	if err != nil {
		log.Fatal(fmt.Errorf("error occured while initializing database object, %v", err))
	}

	//initialize service dependencies
	services := app.NewServices(sqlDB)

	//initialize router
	router := app.NewRouter(services)

	//declared HTTPPort Config as a constant and used here, instead of loading it from environment variable
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", HTTPPort),
		Handler: router,
	}

	slog.InfoContext(ctx, "Starting HTTP Server", slog.Int("port", HTTPPort))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(fmt.Errorf("HTTP Server Closed: %v", err))
	}
}

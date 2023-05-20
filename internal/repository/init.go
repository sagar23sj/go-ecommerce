package repository

import (
	"context"
	"log"

	"github.com/asdine/storm/v3"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/logger"
	"go.uber.org/zap"
)

func InitializeDatabase() (db *storm.DB, err error) {
	db, err = storm.Open("my.db")
	if err != nil {
		log.Printf("error occured while creating database connection: %v", err.Error())
		return nil, err
	}

	//migrate database tables
	db.Init(&Order{})
	db.Init(&Product{})
	db.Init(&OrderItem{})

	//seed products in database
	err = seedDatabase(db)
	if err != nil {
		log.Printf("error occured while seeding products database: %v", err.Error())
		return nil, err
	}

	return db, nil
}

func seedDatabase(db *storm.DB) (err error) {

	products := make([]Product, 0)
	err = db.All(&products)

	//database already has some products, so not adding products again
	if len(products) > 0 {
		return
	}

	products = []Product{
		{Name: "Nike Sneaker", Price: 5000.00, Category: "Premium", Quantity: 20},
		{Name: "Puma Hoodie", Price: 3000.00, Category: "Premium", Quantity: 20},
		{Name: "G-Shock Watch", Price: 8000.00, Category: "Premium", Quantity: 20},
		{Name: "X-Box 360", Price: 25000.00, Category: "Premium", Quantity: 20},
		{Name: "Samsung Smart Watch", Price: 10000.00, Category: "Premium", Quantity: 20},
		{Name: "H&M Sweat Shirt", Price: 1500.00, Category: "Regular", Quantity: 20},
		{Name: "RedTape Sneakers", Price: 1800.00, Category: "Regular", Quantity: 20},
		{Name: "Shirt", Price: 800.00, Category: "Budget", Quantity: 20},
		{Name: "Pant", Price: 1000.00, Category: "Budget", Quantity: 20},
	}

	for _, product := range products {
		err = db.Save(product)
		if err != nil {
			logger.Errorw(context.Background(), "error occured while seeding product in database",
				zap.Error(err),
				zap.String("product_name", product.Name),
			)
		}
	}

	return err
}

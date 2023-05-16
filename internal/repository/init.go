package repository

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDatabase() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Printf("error occured while creating database connection: %v", err.Error())
		return nil, err
	}

	//migrate database tables
	db.AutoMigrate(&Order{}, &Product{}, &OrderItem{})

	//seed products in database
	err = seedDatabase(db)
	if err != nil {
		log.Printf("error occured while seeding products database: %v", err.Error())
		return nil, err
	}

	return db, nil
}

func seedDatabase(db *gorm.DB) (err error) {

	products := make([]Product, 0)
	rowCount := db.Find(&products).RowsAffected

	//database already has some products, so not adding products again
	if rowCount > 0 {
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

	err = db.Create(products).Error
	return err
}

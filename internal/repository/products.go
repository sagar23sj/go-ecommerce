package repository

import (
	"context"
	"time"
)

type ProductStorer interface {
	RepositoryTransaction

	GetProductByID(ctx context.Context, tx Transaction, productID int64) (Product, error)
	ListProducts(ctx context.Context, tx Transaction) ([]Product, error)
	UpdateProductQuantity(ctx context.Context, tx Transaction, productsQuantityMap map[int64]int64) error
}

type Product struct {
	ID        uint `storm:"id,increment"`
	Name      string
	Price     float64
	Category  string
	Quantity  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

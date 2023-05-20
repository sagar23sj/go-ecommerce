package repository

import (
	"context"
	"time"

	"github.com/asdine/storm/v3"
)

type productStore struct {
	BaseRepository
}

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

func NewProductRepo(db *storm.DB) ProductStorer {
	return &productStore{
		BaseRepository: BaseRepository{db},
	}
}

func (ps *productStore) GetProductByID(ctx context.Context, tx Transaction, productID int64) (Product, error) {
	var product Product

	queryExecutor := ps.initiateQueryExecutor(tx)
	err := queryExecutor.One("ID", productID, &product)
	if err != nil && err != storm.ErrNotFound {
		return Product{}, err
	}

	return product, nil
}

func (ps *productStore) ListProducts(ctx context.Context, tx Transaction) ([]Product, error) {
	productList := make([]Product, 0)

	queryExecutor := ps.initiateQueryExecutor(tx)
	err := queryExecutor.All(&productList)
	if err != nil {
		return productList, err
	}

	return productList, nil
}

func (ps *productStore) UpdateProductQuantity(ctx context.Context, tx Transaction, productsQuantityMap map[int64]int64) error {
	queryExecutor := ps.initiateQueryExecutor(tx)

	// Iterate over the map to set the quantity for each product ID
	for productID, quantity := range productsQuantityMap {
		// Update the records with the given product ID
		err := queryExecutor.UpdateField(&Product{ID: uint(productID)}, "Quantity", quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

package repository

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type productStore struct {
	BaseRepository
}

func NewProductRepo(db *storm.DB) repository.ProductStorer {
	return &productStore{
		BaseRepository: BaseRepository{db},
	}
}

func (ps *productStore) GetProductByID(ctx context.Context, tx repository.Transaction, productID int64) (repository.Product, error) {
	var product repository.Product

	queryExecutor := ps.initiateQueryExecutor(tx)
	err := queryExecutor.One("ID", productID, &product)
	if err != nil && err != storm.ErrNotFound {
		return repository.Product{}, err
	}

	return product, nil
}

func (ps *productStore) ListProducts(ctx context.Context, tx repository.Transaction) ([]repository.Product, error) {
	productList := make([]repository.Product, 0)

	queryExecutor := ps.initiateQueryExecutor(tx)
	err := queryExecutor.All(&productList)
	if err != nil {
		return productList, err
	}

	return productList, nil
}

func (ps *productStore) UpdateProductQuantity(ctx context.Context, tx repository.Transaction, productsQuantityMap map[int64]int64) error {
	queryExecutor := ps.initiateQueryExecutor(tx)

	// Iterate over the map to set the quantity for each product ID
	for productID, quantity := range productsQuantityMap {
		// Update the records with the given product ID
		err := queryExecutor.UpdateField(&repository.Product{ID: uint(productID)}, "Quantity", quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

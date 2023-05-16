package repository

import (
	"context"

	"gorm.io/gorm"
)

type productStore struct {
	BaseRepository
}

type ProductStorer interface {
	RepositoryTransaction

	GetProductByID(ctx context.Context, tx *gorm.DB, productID int64) (Product, error)
	ListProducts(ctx context.Context, tx *gorm.DB) ([]Product, error)
	UpdateProductQuantity(ctx context.Context, tx *gorm.DB, productsQuantityMap map[int64]int64) error
}

type Product struct {
	gorm.Model
	Name     string
	Price    float64
	Category string
	Quantity int64
}

func NewProductRepo(db *gorm.DB) ProductStorer {
	return &productStore{
		BaseRepository: BaseRepository{db},
	}
}

func (ps *productStore) GetProductByID(ctx context.Context, tx *gorm.DB, productID int64) (Product, error) {
	var product Product

	queryExecutor := ps.initiateQueryExecutor(tx)
	err := queryExecutor.First(&product, productID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return Product{}, err
	}

	return product, nil
}

func (ps *productStore) ListProducts(ctx context.Context, tx *gorm.DB) ([]Product, error) {
	productList := make([]Product, 0)

	queryExecutor := ps.initiateQueryExecutor(tx)
	err := queryExecutor.Find(&productList).Error
	if err != nil {
		return productList, err
	}

	return productList, nil
}

func (ps *productStore) UpdateProductQuantity(ctx context.Context, tx *gorm.DB, productsQuantityMap map[int64]int64) error {
	queryExecutor := ps.initiateQueryExecutor(tx)

	// Iterate over the map to set the quantity for each product ID
	for productID, quantity := range productsQuantityMap {
		// Update the records with the given product ID
		err := queryExecutor.Model(&Product{}).Where("id = ?", productID).Update("quantity", quantity).Error
		if err != nil {
			return err
		}
	}

	return nil
}

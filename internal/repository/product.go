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

	// queryExecutor := ps.initiateQueryExecutor(tx)
	// err := queryExecutor.First(&product, productID).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return Product{}, err
	// }

	return product, nil
}

func (ps *productStore) ListProducts(ctx context.Context, tx Transaction) ([]Product, error) {
	productList := make([]Product, 0)

	// queryExecutor := ps.initiateQueryExecutor(tx)
	// err := queryExecutor.Find(&productList).Error
	// if err != nil {
	// 	return productList, err
	// }

	return productList, nil
}

func (ps *productStore) UpdateProductQuantity(ctx context.Context, tx Transaction, productsQuantityMap map[int64]int64) error {
	// queryExecutor := ps.initiateQueryExecutor(tx)

	// // Iterate over the map to set the quantity for each product ID
	// for productID, quantity := range productsQuantityMap {
	// 	// Update the records with the given product ID
	// 	err := queryExecutor.Model(&Product{}).Where("id = ?", productID).Update("quantity", quantity).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

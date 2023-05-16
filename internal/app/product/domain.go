package product

import (
	"time"

	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type Product struct {
	ID        int64
	Name      string
	Price     float64
	Category  string
	Quantity  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductList struct {
	Products     []Product
	ProductCount int64
}

func MapRepoObjectToService(repoObj repository.Product) Product {
	return Product{
		ID:        int64(repoObj.ID),
		Name:      repoObj.Name,
		Price:     repoObj.Price,
		Category:  repoObj.Category,
		Quantity:  repoObj.Quantity,
		CreatedAt: repoObj.CreatedAt,
		UpdatedAt: repoObj.CreatedAt,
	}
}

func MapServiceObjectToRepo(product Product) repository.Product {
	return repository.Product{
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
		Quantity: product.Quantity,
	}
}

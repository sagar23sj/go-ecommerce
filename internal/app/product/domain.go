package product

import (
	"github.com/sagar23sj/go-ecommerce/internal/model"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type ProductType string

const (
	PremiumProduct ProductType = "Premium"
	RegularProduct ProductType = "Regular"
	BudgetProduct  ProductType = "Budget"
)

func MapRepoObjectToDto(repoObj repository.Product) model.Product {
	return model.Product{
		ID:        int64(repoObj.ID),
		Name:      repoObj.Name,
		Price:     repoObj.Price,
		Category:  repoObj.Category,
		Quantity:  repoObj.Quantity,
		CreatedAt: repoObj.CreatedAt,
		UpdatedAt: repoObj.CreatedAt,
	}
}

func MapDtoObjectToRepo(product model.Product) repository.Product {
	return repository.Product{
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
		Quantity: product.Quantity,
	}
}

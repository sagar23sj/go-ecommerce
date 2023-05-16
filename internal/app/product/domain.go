package product

import (
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type ProductType string

const (
	PremiumProduct ProductType = "premium"
	RegularProduct ProductType = "regular"
	BudgetProduct  ProductType = "budget"
)

const (
	MaxProductQuantity         = 10
	PremiumProductsForDiscount = 3
)

var ProductTypeMap = map[ProductType]struct{}{
	PremiumProduct: {},
	RegularProduct: {},
	BudgetProduct:  {},
}

func MapRepoObjectToDto(repoObj repository.Product) dto.Product {
	return dto.Product{
		ID:        int64(repoObj.ID),
		Name:      repoObj.Name,
		Price:     repoObj.Price,
		Category:  repoObj.Category,
		Quantity:  repoObj.Quantity,
		CreatedAt: repoObj.CreatedAt,
		UpdatedAt: repoObj.CreatedAt,
	}
}

func MapDtoObjectToRepo(product dto.Product) repository.Product {
	return repository.Product{
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
		Quantity: product.Quantity,
	}
}
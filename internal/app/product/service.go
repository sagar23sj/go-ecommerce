package product

import (
	"context"

	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type service struct {
	productRepo repository.ProductStorer
}

type Service interface {
	GetProductByID(ctx context.Context, productID int64) (Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
}

func NewService(productRepo repository.ProductStorer) Service {
	return &service{
		productRepo: productRepo,
	}
}

func (ps *service) GetProductByID(ctx context.Context, productID int64) (Product, error) {
	productInfoDB, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		return Product{}, nil
	}

	if productInfoDB.ID == 0 {
		return Product{}, apperrors.ProductNotFound{ID: productID}
	}

	productInfo := MapRepoObjectToService(productInfoDB)
	return productInfo, nil
}

func (ps *service) ListProducts(ctx context.Context) ([]Product, error) {
	products := make([]Product, 0)

	productsListDB, err := ps.productRepo.ListProducts(ctx, nil)
	if err != nil {
		return products, err
	}

	for _, productInfo := range productsListDB {
		products = append(products, MapRepoObjectToService(productInfo))
	}

	return products, nil
}

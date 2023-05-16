package api

import (
	"time"

	"github.com/sagar23sj/go-ecommerce/internal/app/product"
)

type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Category  string    `json:"category"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductList struct {
	Products []Product `json:"products"`
}

func MapProductDtoToResponse(productInfo product.Product) Product {
	return Product{
		ID:        productInfo.ID,
		Name:      productInfo.Name,
		Price:     productInfo.Price,
		Category:  productInfo.Category,
		Quantity:  productInfo.Quantity,
		CreatedAt: productInfo.CreatedAt,
		UpdatedAt: productInfo.UpdatedAt,
	}
}

func MapProductListToResponse(products []product.Product) ProductList {
	productList := make([]Product, 0)

	for _, productInfo := range products {
		productList = append(productList, Product{
			ID:        productInfo.ID,
			Name:      productInfo.Name,
			Price:     productInfo.Price,
			Category:  productInfo.Category,
			Quantity:  productInfo.Quantity,
			CreatedAt: productInfo.CreatedAt,
			UpdatedAt: productInfo.UpdatedAt,
		})
	}

	return ProductList{productList}
}

package order

import "time"

type Order struct {
	ID                 int64
	ProductDetails     []ProductInfo
	Amount             float64
	DiscountPercentage int64
	DiscountedAmount   float64
	Status             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ProductInfo struct {
	ProductID int64
	Quantity  int64
}

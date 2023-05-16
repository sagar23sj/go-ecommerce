package dto

import "time"

type Product struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Category  string    `json:"category,omitempty"`
	Quantity  int64     `json:"quantity,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type ProductList struct {
	Products []Product `json:"products"`
}

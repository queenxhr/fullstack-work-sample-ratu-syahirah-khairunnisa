package entity

import "time"

type Product struct {
	ID          int64     `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"` // "IN_STOCK", "LOW_STOCK", "OUT_OF_STOCK"
	Version     int64     `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) ComputeStatus() string {
	if p.Stock <= 0 {
		return "OUT_OF_STOCK"
	}
	if p.Stock <= 5 {
		return "LOW_STOCK"
	}
	return "IN_STOCK"
}

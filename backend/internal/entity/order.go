package entity

import "time"

type Order struct {
	ID             int64     `json:"id"`
	OrderNumber    string    `json:"order_number"`
	ProductID      int64     `json:"product_id"`
	ProductName    string    `json:"product_name"`
	Quantity       int       `json:"quantity"`
	UnitPrice      float64   `json:"unit_price"`
	TotalAmount    float64   `json:"total_amount"`
	UserID         string    `json:"user_id"`
	Status         string    `json:"status"` // "PENDING", "CONFIRMED", "FAILED"
	HandlingMethod string    `json:"handling_method"` // "ATOMIC_UPDATE" or "PESSIMISTIC_LOCK"
	CreatedAt      time.Time `json:"created_at"`
}

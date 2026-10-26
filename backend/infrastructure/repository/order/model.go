package order

import (
	"time"

	"fullstack-backend/internal/entity"
)

type OrderModel struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber    string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_number"`
	ProductID      int64     `gorm:"not null;index" json:"product_id"`
	ProductName    string    `gorm:"type:varchar(255);not null" json:"product_name"`
	Quantity       int       `gorm:"not null;default:1" json:"quantity"`
	UnitPrice      float64   `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	TotalAmount    float64   `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	UserID         string    `gorm:"type:varchar(64);not null" json:"user_id"`
	Status         string    `gorm:"type:varchar(32);not null;default:'CONFIRMED'" json:"status"`
	HandlingMethod string    `gorm:"type:varchar(32);not null" json:"handling_method"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (OrderModel) TableName() string {
	return "orders"
}

func (m OrderModel) ToEntity() *entity.Order {
	return &entity.Order{
		ID:             m.ID,
		OrderNumber:    m.OrderNumber,
		ProductID:      m.ProductID,
		ProductName:    m.ProductName,
		Quantity:       m.Quantity,
		UnitPrice:      m.UnitPrice,
		TotalAmount:    m.TotalAmount,
		UserID:         m.UserID,
		Status:         m.Status,
		HandlingMethod: m.HandlingMethod,
		CreatedAt:      m.CreatedAt,
	}
}

func (OrderModel) FromEntity(e *entity.Order) *OrderModel {
	if e == nil {
		return nil
	}
	return &OrderModel{
		ID:             e.ID,
		OrderNumber:    e.OrderNumber,
		ProductID:      e.ProductID,
		ProductName:    e.ProductName,
		Quantity:       e.Quantity,
		UnitPrice:      e.UnitPrice,
		TotalAmount:    e.TotalAmount,
		UserID:         e.UserID,
		Status:         e.Status,
		HandlingMethod: e.HandlingMethod,
		CreatedAt:      e.CreatedAt,
	}
}

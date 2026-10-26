package product

import (
	"time"

	"fullstack-backend/internal/entity"
)

type ProductModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SKU       string    `gorm:"type:varchar(50);uniqueIndex" json:"sku"`
	Name      string    `gorm:"type:varchar(255);not null;index" json:"name"`
	Category  string    `gorm:"type:varchar(100);not null;index" json:"category"`
	Price     float64   `gorm:"type:numeric(12,2);not null;default:0.00" json:"price"`
	Stock     int       `gorm:"not null;default:0;index" json:"stock"`
	Status    string    `gorm:"type:varchar(50);not null;default:'IN_STOCK'" json:"status"`
	Version   int64     `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ProductModel) TableName() string {
	return "products"
}

func (m ProductModel) ToEntity() *entity.Product {
	return &entity.Product{
		ID:        m.ID,
		SKU:       m.SKU,
		Name:      m.Name,
		Category:  m.Category,
		Price:     m.Price,
		Stock:     m.Stock,
		Status:    m.Status,
		Version:   m.Version,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (ProductModel) FromEntity(e *entity.Product) *ProductModel {
	if e == nil {
		return nil
	}
	return &ProductModel{
		ID:        e.ID,
		SKU:       e.SKU,
		Name:      e.Name,
		Category:  e.Category,
		Price:     e.Price,
		Stock:     e.Stock,
		Status:    e.Status,
		Version:   e.Version,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type listProductModel []ProductModel

func (list listProductModel) ToEntities() []*entity.Product {
	entities := make([]*entity.Product, len(list))
	for i, m := range list {
		entities[i] = m.ToEntity()
	}
	return entities
}

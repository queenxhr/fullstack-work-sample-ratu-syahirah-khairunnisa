package product

import (
	"time"

	"fullstack-backend/internal/entity"
)

type ProductObject struct {
	ID        int64     `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Status    string    `json:"status"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductObject) ParseFromEntity(e *entity.Product) *ProductObject {
	if e == nil {
		return nil
	}
	return &ProductObject{
		ID:        e.ID,
		SKU:       e.SKU,
		Name:      e.Name,
		Category:  e.Category,
		Price:     e.Price,
		Stock:     e.Stock,
		Status:    e.Status,
		Version:   e.Version,
		UpdatedAt: e.UpdatedAt,
	}
}

type ProductListResponse struct {
	Total    int64            `json:"total"`
	Count    int              `json:"count"`
	Products []*ProductObject `json:"products"`
}

func (ProductListResponse) ParseFromEntities(entities []*entity.Product, total int64) *ProductListResponse {
	list := make([]*ProductObject, len(entities))
	for i, e := range entities {
		list[i] = ProductObject{}.ParseFromEntity(e)
	}
	return &ProductListResponse{
		Total:    total,
		Count:    len(list),
		Products: list,
	}
}

package repository

import (
	"context"
	"fullstack-backend/internal/entity"
)

type ProductRepository interface {
	FetchAll(ctx context.Context, limit, offset int, search string) ([]*entity.Product, int64, error)
	FindById(ctx context.Context, id int64) (*entity.Product, error)
	UpdateStock(ctx context.Context, id int64, newStock int) (*entity.Product, error)
	Count(ctx context.Context) (int64, error)
	BatchSeed(ctx context.Context, products []*entity.Product) error
}

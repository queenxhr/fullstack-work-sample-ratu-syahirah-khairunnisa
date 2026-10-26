package repository

import (
	"context"
	"fullstack-backend/internal/entity"
)

type OrderRepository interface {
	CreateWithAtomicStock(ctx context.Context, order *entity.Order) (*entity.Order, error)
	CreateWithPessimisticLock(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetStock(ctx context.Context, productID int64) (int, error)
	ResetStock(ctx context.Context, productID int64, stock int) error
	FetchOrders(ctx context.Context, limit int) ([]*entity.Order, error)
}

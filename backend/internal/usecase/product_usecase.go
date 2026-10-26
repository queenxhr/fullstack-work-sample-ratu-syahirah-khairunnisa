package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/repository"
)

type SSEClient chan string

type IProductUseCase interface {
	FetchAll(ctx context.Context, limit, offset int, search string) ([]*entity.Product, int64, error)
	FindById(ctx context.Context, id int64) (*entity.Product, error)
	UpdateStock(ctx context.Context, id int64, newStock int) (*entity.Product, error)
	SubscribeSSE() SSEClient
	UnsubscribeSSE(client SSEClient)
	BroadcastStockChange(product *entity.Product)
}

type ProductUseCase struct {
	productRepo repository.ProductRepository
	sseClients  map[SSEClient]bool
	sseMu       sync.RWMutex
}

func NewProductUseCase(productRepo repository.ProductRepository) IProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		sseClients:  make(map[SSEClient]bool),
	}
}

func (uc *ProductUseCase) FetchAll(ctx context.Context, limit, offset int, search string) ([]*entity.Product, int64, error) {
	return uc.productRepo.FetchAll(ctx, limit, offset, search)
}

func (uc *ProductUseCase) FindById(ctx context.Context, id int64) (*entity.Product, error) {
	prod, err := uc.productRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if prod == nil {
		return nil, appErrors.New(appErrors.ErrNotFound, "Product not found")
	}
	return prod, nil
}

func (uc *ProductUseCase) UpdateStock(ctx context.Context, id int64, newStock int) (*entity.Product, error) {
	if newStock < 0 {
		return nil, appErrors.New(appErrors.ErrBadRequest, "Stock cannot be negative")
	}

	prod, err := uc.productRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if prod == nil {
		return nil, appErrors.New(appErrors.ErrNotFound, "Product not found")
	}

	updated, err := uc.productRepo.UpdateStock(ctx, id, newStock)
	if err != nil {
		return nil, err
	}

	// Broadcast stock change in real-time to all connected React clients via SSE!
	uc.BroadcastStockChange(updated)

	return updated, nil
}

func (uc *ProductUseCase) SubscribeSSE() SSEClient {
	uc.sseMu.Lock()
	defer uc.sseMu.Unlock()

	client := make(SSEClient, 10)
	uc.sseClients[client] = true
	return client
}

func (uc *ProductUseCase) UnsubscribeSSE(client SSEClient) {
	uc.sseMu.Lock()
	defer uc.sseMu.Unlock()

	if _, ok := uc.sseClients[client]; ok {
		delete(uc.sseClients, client)
		close(client)
	}
}

func (uc *ProductUseCase) BroadcastStockChange(product *entity.Product) {
	uc.sseMu.RLock()
	defer uc.sseMu.RUnlock()

	payload, err := json.Marshal(map[string]interface{}{
		"type":    "STOCK_UPDATED",
		"product": product,
	})
	if err != nil {
		return
	}

	msg := fmt.Sprintf("data: %s\n\n", string(payload))
	for client := range uc.sseClients {
		select {
		case client <- msg:
		default:
			// Non-blocking drop if client buffer is full
		}
	}
}

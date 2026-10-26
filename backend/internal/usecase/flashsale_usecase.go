package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/repository"
)

type BenchmarkResult struct {
	TotalRequests      int           `json:"total_requests"`
	SuccessfulOrders   int           `json:"successful_orders"`
	RejectedOutStock   int           `json:"rejected_out_of_stock"`
	FailedWithErrors   int           `json:"failed_with_errors"`
	InitialStock       int           `json:"initial_stock"`
	RemainingStock     int           `json:"remaining_stock"`
	DurationMs         int64         `json:"duration_ms"`
	RaceConditionOccurred bool       `json:"race_condition_occurred"`
	MethodUsed         string        `json:"method_used"`
	SuccessfulUserIDs  []string      `json:"successful_user_ids"`
}

type IFlashSaleUseCase interface {
	BuyProduct(ctx context.Context, productID int64, userID string, quantity int, method string) (*entity.Order, error)
	ResetStock(ctx context.Context, productID int64, stock int) error
	GetStock(ctx context.Context, productID int64) (int, error)
	Run100ConcurrentBenchmark(ctx context.Context, productID int64, method string) (*BenchmarkResult, error)
	FetchRecentOrders(ctx context.Context, limit int) ([]*entity.Order, error)
}

type FlashSaleUseCase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewFlashSaleUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) IFlashSaleUseCase {
	return &FlashSaleUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (uc *FlashSaleUseCase) BuyProduct(ctx context.Context, productID int64, userID string, quantity int, method string) (*entity.Order, error) {
	if quantity <= 0 {
		return nil, appErrors.New(appErrors.ErrBadRequest, "Quantity must be greater than 0")
	}

	order := &entity.Order{
		OrderNumber:    fmt.Sprintf("ORD-%d-%s", time.Now().UnixNano(), userID),
		ProductID:      productID,
		Quantity:       quantity,
		UserID:         userID,
		Status:         "CONFIRMED",
		HandlingMethod: method,
		CreatedAt:      time.Now(),
	}

	if method == "PESSIMISTIC_LOCK" {
		return uc.orderRepo.CreateWithPessimisticLock(ctx, order)
	}

	// Default to Atomic Update
	order.HandlingMethod = "ATOMIC_UPDATE"
	return uc.orderRepo.CreateWithAtomicStock(ctx, order)
}

func (uc *FlashSaleUseCase) ResetStock(ctx context.Context, productID int64, stock int) error {
	return uc.orderRepo.ResetStock(ctx, productID, stock)
}

func (uc *FlashSaleUseCase) GetStock(ctx context.Context, productID int64) (int, error) {
	return uc.orderRepo.GetStock(ctx, productID)
}

func (uc *FlashSaleUseCase) FetchRecentOrders(ctx context.Context, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.FetchOrders(ctx, limit)
}

// Run100ConcurrentBenchmark fires 100 concurrent requests against a product with 1 remaining stock
func (uc *FlashSaleUseCase) Run100ConcurrentBenchmark(ctx context.Context, productID int64, method string) (*BenchmarkResult, error) {
	initialStock := 1
	// 1. Reset stock to 1
	if err := uc.ResetStock(ctx, productID, initialStock); err != nil {
		return nil, err
	}

	startTime := time.Now()
	totalConcurrent := 100

	var wg sync.WaitGroup
	var mu sync.Mutex

	successCount := 0
	rejectedCount := 0
	errorCount := 0
	var successfulUsers []string

	wg.Add(totalConcurrent)
	for i := 1; i <= totalConcurrent; i++ {
		userID := fmt.Sprintf("buyer-usr-%03d", i)
		go func(uid string) {
			defer wg.Done()

			_, err := uc.BuyProduct(ctx, productID, uid, 1, method)
			mu.Lock()
			defer mu.Unlock()

			if err == nil {
				successCount++
				successfulUsers = append(successfulUsers, uid)
			} else {
				if appErr, ok := err.(*appErrors.DynamicError); ok && appErr.StatusCode == 409 {
					rejectedCount++
				} else {
					errorCount++
				}
			}
		}(userID)
	}

	wg.Wait()
	duration := time.Since(startTime)

	remainingStock, _ := uc.GetStock(ctx, productID)

	// Race condition occurred if more than 1 user succeeded when initial stock was 1, or stock became negative!
	raceConditionOccurred := successCount > initialStock || remainingStock < 0

	result := &BenchmarkResult{
		TotalRequests:         totalConcurrent,
		SuccessfulOrders:      successCount,
		RejectedOutStock:      rejectedCount,
		FailedWithErrors:      errorCount,
		InitialStock:          initialStock,
		RemainingStock:        remainingStock,
		DurationMs:            duration.Milliseconds(),
		RaceConditionOccurred: raceConditionOccurred,
		MethodUsed:            method,
		SuccessfulUserIDs:     successfulUsers,
	}

	return result, nil
}

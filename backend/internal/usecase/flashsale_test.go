package usecase_test

import (
	"context"
	"testing"
	"time"

	orderRepo "fullstack-backend/infrastructure/repository/order"
	productRepo "fullstack-backend/infrastructure/repository/product"
	"fullstack-backend/internal/usecase"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	err = db.AutoMigrate(
		&productRepo.ProductModel{},
		&orderRepo.OrderModel{},
	)
	if err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}

	return db
}

func TestFlashSaleUseCase_100_Concurrent_Buy_1_Stock(t *testing.T) {
	db := setupTestDB(t)

	prodRepo := productRepo.NewPostgresRepository(db)
	ordRepo := orderRepo.NewPostgresRepository(db)
	flashSaleUc := usecase.NewFlashSaleUseCase(ordRepo, prodRepo)

	ctx := context.Background()

	// Seed 1 product with stock = 1
	prod := &productRepo.ProductModel{
		ID:        1,
		SKU:       "SKU-TEST-001",
		Name:      "Limited Flash Sale Item",
		Category:  "Gaming",
		Price:     1000000,
		Stock:     1,
		Status:    "IN_STOCK",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(prod)

	// Run benchmark: 100 concurrent users buying 1 stock
	res, err := flashSaleUc.Run100ConcurrentBenchmark(ctx, 1, "ATOMIC_UPDATE")
	if err != nil {
		t.Fatalf("benchmark failed: %v", err)
	}

	t.Logf("Benchmark Result: Total=%d, Success=%d, Rejected=%d, RemainingStock=%d, Duration=%dms",
		res.TotalRequests, res.SuccessfulOrders, res.RejectedOutStock, res.RemainingStock, res.DurationMs)

	if res.SuccessfulOrders != 1 {
		t.Fatalf("EXPECTED EXACTLY 1 SUCCESSFUL ORDER, GOT %d (RACE CONDITION DETECTED!)", res.SuccessfulOrders)
	}

	if res.RejectedOutStock != 99 {
		t.Fatalf("EXPECTED 99 REJECTED ORDERS (OUT OF STOCK), GOT %d", res.RejectedOutStock)
	}

	if res.RemainingStock != 0 {
		t.Fatalf("EXPECTED REMAINING STOCK TO BE 0, GOT %d", res.RemainingStock)
	}

	if res.RaceConditionOccurred {
		t.Fatalf("TEST FAILED: Race condition occurred!")
	}

	t.Log("✅ SUKSES: 100 concurrent requests aman! 1 sukses, 99 ditolak (409 Conflict), sisa stok = 0 (No Overselling).")
}

package usecase_test

import (
	"context"
	"testing"
	"time"

	webhookRepo "fullstack-backend/infrastructure/repository/webhook"
	"fullstack-backend/internal/entity"
	"fullstack-backend/internal/usecase"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestWebhookUseCase_Idempotency_Deduplication(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	db.AutoMigrate(&webhookRepo.WebhookEventModel{})
	repo := webhookRepo.NewPostgresRepository(db)
	webhookUc := usecase.NewWebhookUseCase(repo)

	ctx := context.Background()

	idempotencyKey := "IDEMP-KEY-TEST-9999"
	event := &entity.WebhookEvent{
		EventID:        idempotencyKey,
		PaymentID:      "PAY-TEST-001",
		OrderID:        "ORD-TEST-001",
		Amount:         500000,
		Currency:       "IDR",
		Status:         "SUCCESS",
		RawPayload:     `{"status": "SUCCESS"}`,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 1. First execution: should be processed as new
	res1, err := webhookUc.ProcessWebhook(ctx, idempotencyKey, event)
	if err != nil {
		t.Fatalf("first webhook execution failed: %v", err)
	}
	if res1.IsDuplicate {
		t.Fatalf("expected first execution NOT to be duplicate")
	}
	if res1.Event.ProcessedCount != 1 {
		t.Fatalf("expected processed count = 1, got: %d", res1.Event.ProcessedCount)
	}

	// 2. Second execution with the SAME idempotency key: should be detected as duplicate!
	res2, err := webhookUc.ProcessWebhook(ctx, idempotencyKey, event)
	if err != nil {
		t.Fatalf("second webhook execution failed: %v", err)
	}
	if !res2.IsDuplicate {
		t.Fatalf("expected second execution to be marked as IS_DUPLICATE=true")
	}
	if res2.Event.ProcessedCount != 2 {
		t.Fatalf("expected processed count = 2, got: %d", res2.Event.ProcessedCount)
	}

	t.Logf("✅ SUKSES: Idempotency berhasil mencegah transaksi ganda! ProcessedCount = %d", res2.Event.ProcessedCount)
}

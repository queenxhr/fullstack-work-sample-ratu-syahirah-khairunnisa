package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/repository"
)

type WebhookResult struct {
	Event       *entity.WebhookEvent `json:"event"`
	IsDuplicate bool                 `json:"is_duplicate"`
	Message     string               `json:"message"`
}

type IWebhookUseCase interface {
	ProcessWebhook(ctx context.Context, idempotencyKey string, event *entity.WebhookEvent) (*WebhookResult, error)
	ListEvents(ctx context.Context, limit int) ([]*entity.WebhookEvent, error)
	ReconcilePendingEvents(ctx context.Context) ([]*entity.WebhookEvent, error)
	SubscribePaymentSSE(paymentID string) SSEClient
	UnsubscribePaymentSSE(paymentID string, client SSEClient)
	BroadcastPaymentStatus(paymentID, status string)
}

type WebhookUseCase struct {
	webhookRepo   repository.WebhookRepository
	paymentSSE    map[string]map[SSEClient]bool
	paymentSSEMu  sync.RWMutex
}

func NewWebhookUseCase(webhookRepo repository.WebhookRepository) IWebhookUseCase {
	return &WebhookUseCase{
		webhookRepo:  webhookRepo,
		paymentSSE:   make(map[string]map[SSEClient]bool),
	}
}

func (uc *WebhookUseCase) ProcessWebhook(ctx context.Context, idempotencyKey string, event *entity.WebhookEvent) (*WebhookResult, error) {
	if idempotencyKey == "" {
		idempotencyKey = event.EventID
	}
	if idempotencyKey == "" {
		return nil, appErrors.New(appErrors.ErrBadRequest, "Idempotency-Key or event_id is required")
	}

	event.IdempotencyKey = idempotencyKey

	// 1. Idempotency Check: Query by idempotency key
	existing, err := uc.webhookRepo.FindEventByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, err
	}

	// 2. If already processed, DO NOT execute side-effects again (idempotent response)
	if existing != nil {
		existing.ProcessedCount++
		updated, _ := uc.webhookRepo.UpdateEvent(ctx, existing)
		return &WebhookResult{
			Event:       updated,
			IsDuplicate: true,
			Message:     fmt.Sprintf("Idempotency Active: Webhook dengan ID '%s' telah diproses sebelumnya (panggilan ke-%d). Side-effects (kredit saldo/pengiriman) tidak dieksekusi ulang.", idempotencyKey, updated.ProcessedCount),
		}, nil
	}

	// 3. New event: process transaction and save
	if event.Status == "" {
		event.Status = "SUCCESS"
	}
	event.ProcessedCount = 1
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	saved, err := uc.webhookRepo.SaveEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	// 4. Notify any frontend listener via Server-Sent Events (SSE)
	uc.BroadcastPaymentStatus(saved.PaymentID, saved.Status)

	return &WebhookResult{
		Event:       saved,
		IsDuplicate: false,
		Message:     "Webhook baru berhasil diproses dan disimpan secara aman.",
	}, nil
}

func (uc *WebhookUseCase) ListEvents(ctx context.Context, limit int) ([]*entity.WebhookEvent, error) {
	return uc.webhookRepo.ListEvents(ctx, limit)
}

// ReconcilePendingEvents simulates the automated Fallback Cron Job / Dead Letter Queue worker
func (uc *WebhookUseCase) ReconcilePendingEvents(ctx context.Context) ([]*entity.WebhookEvent, error) {
	// Query events that remained PENDING for > 1 minute (simulating lost webhooks)
	pending, err := uc.webhookRepo.FindPendingEventsOlderThan(ctx, 0)
	if err != nil {
		return nil, err
	}

	var reconciled []*entity.WebhookEvent
	for _, evt := range pending {
		// Outbound active polling to Payment Gateway Status API
		evt.Status = "SUCCESS"
		evt.RawPayload = evt.RawPayload + " [RECONCILED_VIA_FALLBACK_CRON]"
		updated, err := uc.webhookRepo.UpdateEvent(ctx, evt)
		if err == nil {
			reconciled = append(reconciled, updated)
			uc.BroadcastPaymentStatus(evt.PaymentID, "SUCCESS")
		}
	}

	return reconciled, nil
}

func (uc *WebhookUseCase) SubscribePaymentSSE(paymentID string) SSEClient {
	uc.paymentSSEMu.Lock()
	defer uc.paymentSSEMu.Unlock()

	client := make(SSEClient, 10)
	if _, ok := uc.paymentSSE[paymentID]; !ok {
		uc.paymentSSE[paymentID] = make(map[SSEClient]bool)
	}
	uc.paymentSSE[paymentID][client] = true
	return client
}

func (uc *WebhookUseCase) UnsubscribePaymentSSE(paymentID string, client SSEClient) {
	uc.paymentSSEMu.Lock()
	defer uc.paymentSSEMu.Unlock()

	if clients, ok := uc.paymentSSE[paymentID]; ok {
		delete(clients, client)
		close(client)
		if len(clients) == 0 {
			delete(uc.paymentSSE, paymentID)
		}
	}
}

func (uc *WebhookUseCase) BroadcastPaymentStatus(paymentID, status string) {
	uc.paymentSSEMu.RLock()
	defer uc.paymentSSEMu.RUnlock()

	clients, ok := uc.paymentSSE[paymentID]
	if !ok {
		return
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"type":       "PAYMENT_STATUS_UPDATE",
		"payment_id": paymentID,
		"status":     status,
		"timestamp":  time.Now().Format(time.RFC3339),
	})

	msg := fmt.Sprintf("data: %s\n\n", string(payload))
	for client := range clients {
		select {
		case client <- msg:
		default:
		}
	}
}

package webhook

import (
	"time"

	"fullstack-backend/internal/entity"
)

type WebhookEventModel struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	EventID        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"event_id"`
	IdempotencyKey string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"idempotency_key"`
	PaymentID      string    `gorm:"type:varchar(100);not null;index" json:"payment_id"`
	OrderID        string    `gorm:"type:varchar(100);not null;index" json:"order_id"`
	Amount         float64   `gorm:"type:numeric(12,2);not null" json:"amount"`
	Currency       string    `gorm:"type:varchar(10);not null;default:'IDR'" json:"currency"`
	Status         string    `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	RawPayload     string    `gorm:"type:text" json:"raw_payload"`
	ProcessedCount int       `gorm:"not null;default:1" json:"processed_count"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (WebhookEventModel) TableName() string {
	return "webhook_events"
}

func (m WebhookEventModel) ToEntity() *entity.WebhookEvent {
	return &entity.WebhookEvent{
		ID:             m.ID,
		EventID:        m.EventID,
		IdempotencyKey: m.IdempotencyKey,
		PaymentID:      m.PaymentID,
		OrderID:        m.OrderID,
		Amount:         m.Amount,
		Currency:       m.Currency,
		Status:         m.Status,
		RawPayload:     m.RawPayload,
		ProcessedCount: m.ProcessedCount,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (WebhookEventModel) FromEntity(e *entity.WebhookEvent) *WebhookEventModel {
	if e == nil {
		return nil
	}
	return &WebhookEventModel{
		ID:             e.ID,
		EventID:        e.EventID,
		IdempotencyKey: e.IdempotencyKey,
		PaymentID:      e.PaymentID,
		OrderID:        e.OrderID,
		Amount:         e.Amount,
		Currency:       e.Currency,
		Status:         e.Status,
		RawPayload:     e.RawPayload,
		ProcessedCount: e.ProcessedCount,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

package entity

import "time"

type WebhookEvent struct {
	ID             int64     `json:"id"`
	EventID        string    `json:"event_id"`        // Unique event ID from Payment Gateway
	IdempotencyKey string    `json:"idempotency_key"` // Deduplication key
	PaymentID      string    `json:"payment_id"`
	OrderID        string    `json:"order_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"` // "PENDING", "SUCCESS", "FAILED"
	RawPayload     string    `json:"raw_payload"`
	ProcessedCount int       `json:"processed_count"` // Demonstrates deduplication count
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

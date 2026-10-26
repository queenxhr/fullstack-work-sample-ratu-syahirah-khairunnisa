package repository

import (
	"context"
	"fullstack-backend/internal/entity"
)

type WebhookRepository interface {
	FindEventByIdempotencyKey(ctx context.Context, key string) (*entity.WebhookEvent, error)
	FindEventByEventID(ctx context.Context, eventID string) (*entity.WebhookEvent, error)
	SaveEvent(ctx context.Context, event *entity.WebhookEvent) (*entity.WebhookEvent, error)
	UpdateEvent(ctx context.Context, event *entity.WebhookEvent) (*entity.WebhookEvent, error)
	ListEvents(ctx context.Context, limit int) ([]*entity.WebhookEvent, error)
	FindPendingEventsOlderThan(ctx context.Context, minutes int) ([]*entity.WebhookEvent, error)
}

package webhook

import (
	"context"
	"fmt"
	"time"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/repository"

	"gorm.io/gorm"
)

type postgresRepository struct {
	dbConn *gorm.DB
}

func NewPostgresRepository(dbConn *gorm.DB) repository.WebhookRepository {
	return &postgresRepository{
		dbConn: dbConn,
	}
}

func (repo *postgresRepository) FindEventByIdempotencyKey(ctx context.Context, key string) (*entity.WebhookEvent, error) {
	var model WebhookEventModel
	err := repo.dbConn.WithContext(ctx).First(&model, "idempotency_key = ?", key).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to find webhook by idempotency key: %v", err))
	}
	return model.ToEntity(), nil
}

func (repo *postgresRepository) FindEventByEventID(ctx context.Context, eventID string) (*entity.WebhookEvent, error) {
	var model WebhookEventModel
	err := repo.dbConn.WithContext(ctx).First(&model, "event_id = ?", eventID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to find webhook by event id: %v", err))
	}
	return model.ToEntity(), nil
}

func (repo *postgresRepository) SaveEvent(ctx context.Context, event *entity.WebhookEvent) (*entity.WebhookEvent, error) {
	model := WebhookEventModel{}.FromEntity(event)
	err := repo.dbConn.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to save webhook event: %v", err))
	}
	return model.ToEntity(), nil
}

func (repo *postgresRepository) UpdateEvent(ctx context.Context, event *entity.WebhookEvent) (*entity.WebhookEvent, error) {
	model := WebhookEventModel{}.FromEntity(event)
	err := repo.dbConn.WithContext(ctx).Model(&WebhookEventModel{}).
		Where("id = ?", model.ID).
		Updates(map[string]interface{}{
			"status":          model.Status,
			"processed_count": model.ProcessedCount,
			"updated_at":      time.Now(),
		}).Error
	if err != nil {
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to update webhook event: %v", err))
	}
	return model.ToEntity(), nil
}

func (repo *postgresRepository) ListEvents(ctx context.Context, limit int) ([]*entity.WebhookEvent, error) {
	var models []WebhookEventModel
	query := repo.dbConn.WithContext(ctx).Order("id desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*entity.WebhookEvent, len(models))
	for i, m := range models {
		result[i] = m.ToEntity()
	}
	return result, nil
}

func (repo *postgresRepository) FindPendingEventsOlderThan(ctx context.Context, minutes int) ([]*entity.WebhookEvent, error) {
	var models []WebhookEventModel
	threshold := time.Now().Add(-time.Duration(minutes) * time.Minute)
	err := repo.dbConn.WithContext(ctx).
		Where("status = ? AND created_at <= ?", "PENDING", threshold).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.WebhookEvent, len(models))
	for i, m := range models {
		result[i] = m.ToEntity()
	}
	return result, nil
}

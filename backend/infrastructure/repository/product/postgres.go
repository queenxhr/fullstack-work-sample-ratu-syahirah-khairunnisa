package product

import (
	"context"
	"fmt"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/repository"
	"gorm.io/gorm"
)

type postgresRepository struct {
	dbConn *gorm.DB
}

func NewPostgresRepository(dbConn *gorm.DB) repository.ProductRepository {
	return &postgresRepository{
		dbConn: dbConn,
	}
}

func NewMysqlRepository(dbConn *gorm.DB) repository.ProductRepository {
	return NewPostgresRepository(dbConn)
}

func (repo *postgresRepository) FetchAll(ctx context.Context, limit, offset int, search string) ([]*entity.Product, int64, error) {
	var models []ProductModel
	var total int64

	query := repo.dbConn.WithContext(ctx).Model(&ProductModel{})
	if search != "" {
		likePattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR sku LIKE ? OR category LIKE ?", likePattern, likePattern, likePattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to count products: %v", err))
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("id asc").Find(&models).Error
	if err != nil {
		return nil, 0, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to fetch products: %v", err))
	}

	return listProductModel(models).ToEntities(), total, nil
}

func (repo *postgresRepository) FindById(ctx context.Context, id int64) (*entity.Product, error) {
	var model ProductModel
	err := repo.dbConn.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to find product by id: %v", err))
	}
	return model.ToEntity(), nil
}

func (repo *postgresRepository) UpdateStock(ctx context.Context, id int64, newStock int) (*entity.Product, error) {
	status := "IN_STOCK"
	if newStock <= 0 {
		status = "OUT_OF_STOCK"
	} else if newStock <= 5 {
		status = "LOW_STOCK"
	}

	res := repo.dbConn.WithContext(ctx).Model(&ProductModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"stock":      newStock,
			"status":     status,
			"version":    gorm.Expr("version + 1"),
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		})

	if res.Error != nil {
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to update product stock: %v", res.Error))
	}

	return repo.FindById(ctx, id)
}

func (repo *postgresRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := repo.dbConn.WithContext(ctx).Model(&ProductModel{}).Count(&count).Error
	return count, err
}

func (repo *postgresRepository) BatchSeed(ctx context.Context, products []*entity.Product) error {
	models := make([]ProductModel, len(products))
	for i, p := range products {
		models[i] = *ProductModel{}.FromEntity(p)
	}
	// Insert in batches of 500
	return repo.dbConn.WithContext(ctx).CreateInBatches(models, 500).Error
}

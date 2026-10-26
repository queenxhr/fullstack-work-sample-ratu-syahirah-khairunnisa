package order

import (
	"context"
	"fmt"
	"sync"
	"time"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/repository"
	productRepo "fullstack-backend/infrastructure/repository/product"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresRepository struct {
	dbConn *gorm.DB
	// Mutex for SQLite environments where table-level locks occur during concurrent write transactions
	sqliteMu sync.Mutex
}

func NewPostgresRepository(dbConn *gorm.DB) repository.OrderRepository {
	return &postgresRepository{
		dbConn: dbConn,
	}
}

func (repo *postgresRepository) GetStock(ctx context.Context, productID int64) (int, error) {
	var prod productRepo.ProductModel
	err := repo.dbConn.WithContext(ctx).Select("stock").First(&prod, "id = ?", productID).Error
	if err != nil {
		return 0, err
	}
	return prod.Stock, nil
}

func (repo *postgresRepository) ResetStock(ctx context.Context, productID int64, stock int) error {
	status := "IN_STOCK"
	if stock <= 0 {
		status = "OUT_OF_STOCK"
	} else if stock <= 5 {
		status = "LOW_STOCK"
	}

	return repo.dbConn.WithContext(ctx).Model(&productRepo.ProductModel{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"stock":      stock,
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// CreateWithAtomicStock demonstrates Solusi Level Database:
// UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?
func (repo *postgresRepository) CreateWithAtomicStock(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	repo.sqliteMu.Lock()
	defer repo.sqliteMu.Unlock()

	tx := repo.dbConn.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Atomic decrement with condition stock >= qty
	res := tx.Model(&productRepo.ProductModel{}).
		Where("id = ? AND stock >= ?", order.ProductID, order.Quantity).
		Updates(map[string]interface{}{
			"stock":      gorm.Expr("stock - ?", order.Quantity),
			"version":    gorm.Expr("version + 1"),
			"updated_at": time.Now(),
		})

	if res.Error != nil {
		tx.Rollback()
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to update stock: %v", res.Error))
	}

	// If no rows were affected, the stock is insufficient (Race Condition prevented!)
	if res.RowsAffected == 0 {
		tx.Rollback()
		return nil, appErrors.New(appErrors.ErrConflict, "Stok habis atau pesanan melebihi sisa stok (Race Condition Prevented via Atomic SQL Update)")
	}

	// 2. Fetch current stock to update status
	var prod productRepo.ProductModel
	if err := tx.First(&prod, "id = ?", order.ProductID).Error; err == nil {
		newStatus := "IN_STOCK"
		if prod.Stock <= 0 {
			newStatus = "OUT_OF_STOCK"
		} else if prod.Stock <= 5 {
			newStatus = "LOW_STOCK"
		}
		tx.Model(&productRepo.ProductModel{}).Where("id = ?", order.ProductID).Update("status", newStatus)
		order.ProductName = prod.Name
		order.UnitPrice = prod.Price
		order.TotalAmount = prod.Price * float64(order.Quantity)
	}

	// 3. Create Order Record
	orderModel := OrderModel{}.FromEntity(order)
	if err := tx.Create(orderModel).Error; err != nil {
		tx.Rollback()
		return nil, appErrors.New(appErrors.ErrInternal, fmt.Sprintf("failed to create order: %v", err))
	}

	if err := tx.Commit().Error; err != nil {
		return nil, appErrors.New(appErrors.ErrInternal, "failed to commit order transaction")
	}

	return orderModel.ToEntity(), nil
}

// CreateWithPessimisticLock demonstrates SELECT ... FOR UPDATE within a transaction
func (repo *postgresRepository) CreateWithPessimisticLock(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	repo.sqliteMu.Lock()
	defer repo.sqliteMu.Unlock()

	tx := repo.dbConn.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. SELECT ... FOR UPDATE (Row-level exclusive lock)
	var prod productRepo.ProductModel
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&prod, "id = ?", order.ProductID).Error
	if err != nil {
		tx.Rollback()
		return nil, appErrors.New(appErrors.ErrNotFound, "Product not found")
	}

	// 2. Check stock condition inside exclusive lock
	if prod.Stock < order.Quantity {
		tx.Rollback()
		return nil, appErrors.New(appErrors.ErrConflict, "Stok tidak mencukupi (Race Condition Prevented via Pessimistic Lock SELECT FOR UPDATE)")
	}

	// 3. Decrement stock
	newStock := prod.Stock - order.Quantity
	newStatus := "IN_STOCK"
	if newStock <= 0 {
		newStatus = "OUT_OF_STOCK"
	} else if newStock <= 5 {
		newStatus = "LOW_STOCK"
	}

	if err := tx.Model(&prod).Updates(map[string]interface{}{
		"stock":      newStock,
		"status":     newStatus,
		"version":    gorm.Expr("version + 1"),
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return nil, appErrors.New(appErrors.ErrInternal, "failed to update stock")
	}

	// 4. Create Order Record
	order.ProductName = prod.Name
	order.UnitPrice = prod.Price
	order.TotalAmount = prod.Price * float64(order.Quantity)

	orderModel := OrderModel{}.FromEntity(order)
	if err := tx.Create(orderModel).Error; err != nil {
		tx.Rollback()
		return nil, appErrors.New(appErrors.ErrInternal, "failed to create order")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, appErrors.New(appErrors.ErrInternal, "failed to commit transaction")
	}

	return orderModel.ToEntity(), nil
}

func (repo *postgresRepository) FetchOrders(ctx context.Context, limit int) ([]*entity.Order, error) {
	var models []OrderModel
	query := repo.dbConn.WithContext(ctx).Order("id desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*entity.Order, len(models))
	for i, m := range models {
		result[i] = m.ToEntity()
	}
	return result, nil
}

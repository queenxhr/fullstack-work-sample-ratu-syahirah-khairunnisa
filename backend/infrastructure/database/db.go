package database

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"fullstack-backend/config"
	orderRepo "fullstack-backend/infrastructure/repository/order"
	productRepo "fullstack-backend/infrastructure/repository/product"
	webhookRepo "fullstack-backend/infrastructure/repository/webhook"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	if cfg.DBDriver == "postgres" || cfg.DBDriver == "auto" {
		dsn := cfg.GetDSN()
		log.Printf("Connecting to PostgreSQL at %s:%s (database: %s)...", cfg.DBHost, cfg.DBPort, cfg.DBName)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			log.Println("PostgreSQL connection established successfully!")
		} else {
			if cfg.DBDriver == "postgres" {
				return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
			}
			log.Printf("PostgreSQL connection failed (%v). Falling back to embedded SQLite...", err)
		}
	}

	if db == nil {
		sqliteFile := "ecommerce_work_sample.db"
		log.Printf("Initializing SQLite database: %s", sqliteFile)
		db, err = gorm.Open(sqlite.Open(sqliteFile), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SQLite: %w", err)
		}
		log.Println("SQLite database initialized successfully!")
	}

	// Run auto migrations
	if err := AutoMigrate(db); err != nil {
		log.Printf("Warning: AutoMigrate encountered an issue: %v\n", err)
	}

	// Seed 5,000+ realistic products if dataset is empty
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		SeedProductsIfEmpty(ctx, db)
	}()

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&productRepo.ProductModel{},
		&orderRepo.OrderModel{},
		&webhookRepo.WebhookEventModel{},
	)
	if err != nil {
		return fmt.Errorf("automigration failed: %w", err)
	}

	log.Println("Database migration completed successfully!")
	return nil
}

func SeedProductsIfEmpty(ctx context.Context, db *gorm.DB) {
	var count int64
	db.WithContext(ctx).Model(&productRepo.ProductModel{}).Count(&count)
	if count >= 5000 {
		log.Printf("Database already contains %d products. Skipping seeding.", count)
		return
	}

	target := 5000
	needed := target - int(count)
	log.Printf("Seeding %d products to satisfy the 5,000+ dataset requirement (Soal #1)...", needed)

	categories := []string{"Smartphones", "Laptops", "Audio & Wearables", "Gaming & PC", "Accessories", "Home Appliances"}
	brandPrefixes := []string{"Apple", "Sony", "Logitech", "Samsung", "Asus", "Dell", "Bose", "Xiaomi", "Razer", "LG"}
	productTypes := []string{"Pro Max", "Wireless Ultra", "Gaming Edition", "Studio Monitor", "Ergo Mechanical", "OLED Smart", "Thunderbolt Dock", "Noise Cancelling"}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var batch []productRepo.ProductModel
	batchSize := 500

	startIndex := int(count) + 1
	for i := 0; i < needed; i++ {
		idx := startIndex + i
		cat := categories[rng.Intn(len(categories))]
		brand := brandPrefixes[rng.Intn(len(brandPrefixes))]
		ptype := productTypes[rng.Intn(len(productTypes))]
		name := fmt.Sprintf("%s %s %s - Model #%04d", brand, ptype, cat, idx)
		price := float64(rng.Intn(15000000)+150000) / 100.0 // 1500.00 to 150,000.00

		// Special case: Product #1 has stock = 1 for Flash Sale 100-concurrency demonstration!
		stock := rng.Intn(150) + 1
		if idx == 1 {
			stock = 1
			name = "🔥 [FLASH SALE EXCLUSIVE] Sony PlayStation 5 Pro Limited (1 Remaining Unit)"
			price = 8999000.00
			cat = "Gaming & PC"
		}

		status := "IN_STOCK"
		if stock <= 0 {
			status = "OUT_OF_STOCK"
		} else if stock <= 5 {
			status = "LOW_STOCK"
		}

		prod := productRepo.ProductModel{
			SKU:       fmt.Sprintf("SKU-%05d", idx),
			Name:      name,
			Category:  cat,
			Price:     price,
			Stock:     stock,
			Status:    status,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		batch = append(batch, prod)

		if len(batch) >= batchSize || i == needed-1 {
			if err := db.WithContext(ctx).Create(&batch).Error; err != nil {
				log.Printf("Error inserting product batch: %v", err)
			}
			batch = batch[:0]
		}
	}

	var newTotal int64
	db.WithContext(ctx).Model(&productRepo.ProductModel{}).Count(&newTotal)
	log.Printf("Product seeding complete! Total products in catalog: %d", newTotal)
}

package main

import (
	"log"

	"fullstack-backend/config"
	handler "fullstack-backend/handler/v1"
	"fullstack-backend/infrastructure/database"
	orderRepo "fullstack-backend/infrastructure/repository/order"
	productRepo "fullstack-backend/infrastructure/repository/product"
	tokenRepo "fullstack-backend/infrastructure/repository/token"
	webhookRepo "fullstack-backend/infrastructure/repository/webhook"
	"fullstack-backend/internal/usecase"
)

func main() {
	log.Println("================================================================")
	log.Println("🚀 Fullstack Work Sample - Senior React & Go Backend Service")
	log.Println("Candidate: Ratu Syahirah Khairunnisa")
	log.Println("================================================================")

	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database Connection (PostgreSQL with SQLite auto-fallback)
	dbConn, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}

	// 3. Initialize Repositories (Clean Architecture: Dependency Inversion via interfaces)
	productRepository := productRepo.NewPostgresRepository(dbConn)
	orderRepository := orderRepo.NewPostgresRepository(dbConn)
	webhookRepository := webhookRepo.NewPostgresRepository(dbConn)
	tokenRepository := tokenRepo.NewInMemoryRevocationStore()

	// 4. Initialize UseCases (inject semua dependencies)
	productUseCase := usecase.NewProductUseCase(productRepository)
	flashSaleUseCase := usecase.NewFlashSaleUseCase(orderRepository, productRepository)
	authUseCase := usecase.NewAuthUseCase(tokenRepository, cfg)
	webhookUseCase := usecase.NewWebhookUseCase(webhookRepository)

	// 5. Initialize Router & HTTP Handlers
	router := handler.GetRouters(handler.ServiceUseCases{
		ProductUseCase:   productUseCase,
		FlashSaleUseCase: flashSaleUseCase,
		AuthUseCase:      authUseCase,
		WebhookUseCase:   webhookUseCase,
	})

	// 6. Run HTTP Server
	log.Printf("Server is running on port :%s (http://localhost:%s)...", cfg.AppPort, cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}

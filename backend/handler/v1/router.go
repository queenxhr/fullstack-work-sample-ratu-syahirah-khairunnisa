package handler

import (
	"net/http"

	authHttp "fullstack-backend/handler/v1/auth/http"
	flashsaleHttp "fullstack-backend/handler/v1/flashsale/http"
	productHttp "fullstack-backend/handler/v1/product/http"
	webhookHttp "fullstack-backend/handler/v1/webhook/http"
	"fullstack-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ServiceUseCases struct {
	ProductUseCase   usecase.IProductUseCase
	FlashSaleUseCase usecase.IFlashSaleUseCase
	AuthUseCase      usecase.IAuthUseCase
	WebhookUseCase   usecase.IWebhookUseCase
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Idempotency-Key, X-Idempotency-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func GetRouters(uc ServiceUseCases) *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "healthy",
			"service": "Fullstack Work Sample - Senior Fullstack Engineer Assessment",
			"version": "1.0.0",
		})
	})

	apiV1 := r.Group("/api/v1")
	{
		// 1. Products endpoints (Soal #1: 5,000+ items real-time table)
		productHandler := productHttp.NewHandler(uc.ProductUseCase)
		productsGroup := apiV1.Group("/products")
		{
			productsGroup.GET("", productHandler.FetchAll)
			productsGroup.GET("/:id", productHandler.FindById)
			productsGroup.PUT("/:id/stock", productHandler.UpdateStock)
			productsGroup.GET("/stream", productHandler.StreamSSE) // Server-Sent Events live stream
		}

		// 2. Flash Sale Concurrency endpoints (Soal #2: 100 users buy 1 item race condition handling)
		flashSaleHandler := flashsaleHttp.NewHandler(uc.FlashSaleUseCase)
		flashSaleGroup := apiV1.Group("/flashsale")
		{
			flashSaleGroup.POST("/buy", flashSaleHandler.Buy)
			flashSaleGroup.POST("/reset-stock", flashSaleHandler.ResetStock)
			flashSaleGroup.GET("/stock/:id", flashSaleHandler.GetStock)
			flashSaleGroup.POST("/benchmark-100", flashSaleHandler.Run100Benchmark)
			flashSaleGroup.GET("/orders", flashSaleHandler.FetchOrders)
		}

		// 3. Auth endpoints (Soal #2: Stateless JWT + Redis-style JTI Revocation)
		authHandler := authHttp.NewHandler(uc.AuthUseCase)
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/blacklist", authHandler.ListBlacklist)

			// Protected routes
			protected := authGroup.Group("")
			protected.Use(authHttp.JWTMiddleware(uc.AuthUseCase))
			{
				protected.GET("/profile", authHandler.GetProfile)
				protected.POST("/logout", authHandler.Logout)
			}
		}

		// 4. Webhook & Payment Resilience endpoints (Soal #3: Idempotency, DLQ, and SSE)
		webhookHandler := webhookHttp.NewHandler(uc.WebhookUseCase)
		webhookGroup := apiV1.Group("/webhook")
		{
			webhookGroup.POST("/payment", webhookHandler.HandlePaymentWebhook)
			webhookGroup.GET("/events", webhookHandler.ListEvents)
			webhookGroup.POST("/reconcile", webhookHandler.TriggerReconciliation)
			webhookGroup.GET("/stream/:payment_id", webhookHandler.StreamPaymentStatus)
		}
	}

	return r
}

package flashsale

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/usecase"
	"fullstack-backend/utils"
	"github.com/gin-gonic/gin"
)

type IFlashSaleHandler interface {
	Buy(ctx *gin.Context)
	ResetStock(ctx *gin.Context)
	GetStock(ctx *gin.Context)
	Run100Benchmark(ctx *gin.Context)
	FetchOrders(ctx *gin.Context)
}

type FlashSaleHandler struct {
	flashSaleUc usecase.IFlashSaleUseCase
}

func NewHandler(flashSaleUc usecase.IFlashSaleUseCase) IFlashSaleHandler {
	return &FlashSaleHandler{
		flashSaleUc: flashSaleUc,
	}
}

func (h *FlashSaleHandler) Buy(ctx *gin.Context) {
	var reqBody BuyRequestDto
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, err.Error()))
		return
	}

	if reqBody.UserID == "" {
		reqBody.UserID = fmt.Sprintf("buyer-%d", time.Now().UnixNano()%10000)
	}
	if reqBody.Quantity <= 0 {
		reqBody.Quantity = 1
	}

	order, err := h.flashSaleUc.BuyProduct(ctx.Request.Context(), reqBody.ProductID, reqBody.UserID, reqBody.Quantity, reqBody.Method)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusCreated, order, "Pembelian Flash Sale Berhasil! Stok berhasil diproses tanpa race condition.")
}

func (h *FlashSaleHandler) ResetStock(ctx *gin.Context) {
	var reqBody ResetStockRequestDto
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, err.Error()))
		return
	}

	if err := h.flashSaleUc.ResetStock(ctx.Request.Context(), reqBody.ProductID, reqBody.Stock); err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, gin.H{
		"product_id": reqBody.ProductID,
		"new_stock":  reqBody.Stock,
	}, fmt.Sprintf("Stok produk #%d berhasil di-reset menjadi %d untuk pengujian Flash Sale", reqBody.ProductID, reqBody.Stock))
}

func (h *FlashSaleHandler) GetStock(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, "Invalid product ID"))
		return
	}

	stock, err := h.flashSaleUc.GetStock(ctx.Request.Context(), id)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, gin.H{
		"product_id": id,
		"stock":      stock,
	})
}

// Run100Benchmark executes 100 concurrent checkout goroutines against 1 single remaining stock
func (h *FlashSaleHandler) Run100Benchmark(ctx *gin.Context) {
	var reqBody BenchmarkRequestDto
	_ = ctx.ShouldBindJSON(&reqBody)

	productID := reqBody.ProductID
	if productID <= 0 {
		productID = 1 // Default to Flash Sale Exclusive Item #1
	}

	method := reqBody.Method
	if method == "" {
		method = "ATOMIC_UPDATE"
	}

	result, err := h.flashSaleUc.Run100ConcurrentBenchmark(ctx.Request.Context(), productID, method)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, result, "Hasil Pengujian Concurrency 100 User terhadap 1 Sisa Barang")
}

func (h *FlashSaleHandler) FetchOrders(ctx *gin.Context) {
	orders, err := h.flashSaleUc.FetchRecentOrders(ctx.Request.Context(), 50)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, orders)
}

package product

import (
	"io"
	"net/http"
	"strconv"

	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/usecase"
	"fullstack-backend/utils"
	"github.com/gin-gonic/gin"
)

type IProductHandler interface {
	FetchAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	UpdateStock(ctx *gin.Context)
	StreamSSE(ctx *gin.Context)
}

type ProductHandler struct {
	productUc usecase.IProductUseCase
}

func NewHandler(productUseCase usecase.IProductUseCase) IProductHandler {
	return &ProductHandler{
		productUc: productUseCase,
	}
}

func (h *ProductHandler) FetchAll(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "5000") // Default to 5000 for Soal #1 testing!
	offsetStr := ctx.DefaultQuery("offset", "0")
	search := ctx.Query("search")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	products, total, err := h.productUc.FetchAll(ctx.Request.Context(), limit, offset, search)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, ProductListResponse{}.ParseFromEntities(products, total))
}

func (h *ProductHandler) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, "Invalid product ID"))
		return
	}

	product, err := h.productUc.FindById(ctx.Request.Context(), id)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, ProductObject{}.ParseFromEntity(product))
}

func (h *ProductHandler) UpdateStock(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, "Invalid product ID"))
		return
	}

	var reqBody UpdateStockRequestDto
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, err.Error()))
		return
	}

	updated, err := h.productUc.UpdateStock(ctx.Request.Context(), id, reqBody.Stock)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, ProductObject{}.ParseFromEntity(updated), "Stock updated successfully")
}

// StreamSSE streams real-time product updates using Server-Sent Events
func (h *ProductHandler) StreamSSE(ctx *gin.Context) {
	client := h.productUc.SubscribeSSE()
	defer h.productUc.UnsubscribeSSE(client)

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	// Send initial greeting
	ctx.Writer.Write([]byte("event: connected\ndata: {\"message\": \"Connected to real-time product stream\"}\n\n"))
	ctx.Writer.Flush()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client:
			if !ok {
				return false
			}
			w.Write([]byte(msg))
			return true
		case <-ctx.Request.Context().Done():
			return false
		}
	})
}

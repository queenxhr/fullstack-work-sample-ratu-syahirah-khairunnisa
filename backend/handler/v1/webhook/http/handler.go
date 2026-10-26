package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fullstack-backend/internal/entity"
	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/usecase"
	"fullstack-backend/utils"
	"github.com/gin-gonic/gin"
)

type IWebhookHandler interface {
	HandlePaymentWebhook(ctx *gin.Context)
	ListEvents(ctx *gin.Context)
	TriggerReconciliation(ctx *gin.Context)
	StreamPaymentStatus(ctx *gin.Context)
}

type WebhookHandler struct {
	webhookUc usecase.IWebhookUseCase
}

func NewHandler(webhookUc usecase.IWebhookUseCase) IWebhookHandler {
	return &WebhookHandler{
		webhookUc: webhookUc,
	}
}

func (h *WebhookHandler) HandlePaymentWebhook(ctx *gin.Context) {
	// 1. Extract Idempotency-Key from header or query or body
	idempotencyKey := ctx.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		idempotencyKey = ctx.GetHeader("X-Idempotency-Key")
	}

	var reqBody PaymentWebhookRequestDto
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, err.Error()))
		return
	}

	if idempotencyKey == "" {
		idempotencyKey = reqBody.EventID
	}
	if idempotencyKey == "" {
		idempotencyKey = fmt.Sprintf("evt-pg-%s", reqBody.PaymentID)
	}

	raw, _ := json.Marshal(reqBody)
	event := &entity.WebhookEvent{
		EventID:        idempotencyKey,
		PaymentID:      reqBody.PaymentID,
		OrderID:        reqBody.OrderID,
		Amount:         reqBody.Amount,
		Currency:       reqBody.Currency,
		Status:         reqBody.Status,
		RawPayload:     string(raw),
	}

	res, err := h.webhookUc.ProcessWebhook(ctx.Request.Context(), idempotencyKey, event)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	statusCode := http.StatusOK
	if !res.IsDuplicate {
		statusCode = http.StatusCreated
	}

	utils.JSON(ctx, statusCode, res, res.Message)
}

func (h *WebhookHandler) ListEvents(ctx *gin.Context) {
	events, err := h.webhookUc.ListEvents(ctx.Request.Context(), 50)
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, events)
}

func (h *WebhookHandler) TriggerReconciliation(ctx *gin.Context) {
	reconciled, err := h.webhookUc.ReconcilePendingEvents(ctx.Request.Context())
	if err != nil {
		utils.SendJSONErrorResponse(ctx, err)
		return
	}

	utils.JSON(ctx, http.StatusOK, gin.H{
		"reconciled_count": len(reconciled),
		"reconciled_events": reconciled,
		"executed_at": time.Now().Format(time.RFC3339),
	}, "Proses Fallback Cron Job / Reconciliation Selesai. Transaksi menggantung berhasil diverifikasi ulang.")
}

func (h *WebhookHandler) StreamPaymentStatus(ctx *gin.Context) {
	paymentID := ctx.Param("payment_id")
	if paymentID == "" {
		utils.SendJSONErrorResponse(ctx, appErrors.New(appErrors.ErrBadRequest, "payment_id required"))
		return
	}

	client := h.webhookUc.SubscribePaymentSSE(paymentID)
	defer h.webhookUc.UnsubscribePaymentSSE(paymentID, client)

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	ctx.Writer.Write([]byte(fmt.Sprintf("event: connected\ndata: {\"payment_id\": \"%s\", \"listening\": true}\n\n", paymentID)))
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

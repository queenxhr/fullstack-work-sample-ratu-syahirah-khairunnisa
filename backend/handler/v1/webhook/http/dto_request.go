package webhook

type PaymentWebhookRequestDto struct {
	EventID   string  `json:"event_id"`
	PaymentID string  `json:"payment_id" binding:"required"`
	OrderID   string  `json:"order_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"` // "PENDING", "SUCCESS", "FAILED"
}

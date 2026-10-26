package flashsale

type BuyRequestDto struct {
	ProductID int64  `json:"product_id" binding:"required"`
	UserID    string `json:"user_id"`
	Quantity  int    `json:"quantity"`
	Method    string `json:"method"` // "ATOMIC_UPDATE" or "PESSIMISTIC_LOCK"
}

type ResetStockRequestDto struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Stock     int   `json:"stock" binding:"min=0"`
}

type BenchmarkRequestDto struct {
	ProductID int64  `json:"product_id"`
	Method    string `json:"method"` // "ATOMIC_UPDATE" or "PESSIMISTIC_LOCK"
}

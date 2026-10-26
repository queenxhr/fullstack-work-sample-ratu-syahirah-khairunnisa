package product

type UpdateStockRequestDto struct {
	Stock int `json:"stock" binding:"min=0"`
}

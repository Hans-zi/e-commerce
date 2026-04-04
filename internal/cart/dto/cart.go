package dto

type Cart struct {
	ID    string      `json:"id"`
	Lines []*CartLine `json:"lines"`
}
type CartLine struct {
	ProductID string  `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  uint    `json:"quantity"`
}
type AddProductReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  uint   `json:"quantity" binding:"required"`
}

type RemoveProductReq struct {
	ProductID string `json:"product_id" binding:"required"`
}

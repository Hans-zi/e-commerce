package dto

type SeckillReq struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type SeckillOrderAsyncMsg struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

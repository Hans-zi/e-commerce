package dto

import "time"

type Order struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	Total          float64 `json:"total"`
	Status         string  `json:"status"`          // pending, paid, shipped, completed, canceled
	PaymentStatus  string  `json:"payment_status"`  // unpaid, paid, refunded
	ShippingStatus string  `json:"shipping_status"` // unshipped, shipped, received

	Lines []*OrderLine `json:"lines"`
}

type OrderLine struct {
	Quantity int `json:"quantity"`
	Product  Product
}
type PlaceOrderReq struct {
	Lines []PlaceOrderLineReq `json:"lines"`
}

type PlaceOrderLineReq struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type AutoCancelMsg struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	TimeStamp time.Time `json:"time_stamp"`
}

type RefundMsg struct {
	OrderID string `json:"order_id"`
}

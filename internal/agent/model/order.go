package model

import (
	"time"
)

// Order 订单模型
type Order struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Total          float64    `json:"total"`
	Status         string     `json:"status"`          // pending, paid, shipped, completed, canceled
	PaymentStatus  string     `json:"payment_status"`  // unpaid, paid, refunded
	ShippingStatus string     `json:"shipping_status"` // unshipped, shipped, received
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`

	Lines []*OrderLine `json:"lines" gorm:"foreignKey:OrderID"`
}

// OrderLines 订单项
type OrderLine struct {
	ID        string   `json:"id"`
	OrderID   string   `json:"order_id"`
	ProductID string   `json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductID;references:ID"`
	Quantity  int      `json:"quantity"`
	Price     float64  `json:"price"`
}

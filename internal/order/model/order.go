package model

import (
	"e-commerce/internal/consts"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	ID             string     `gorm:"primaryKey" json:"id"`
	UserID         string     `gorm:"type:char(36);not null;index" json:"user_id"`
	Total          float64    `gorm:"not null" json:"total"`
	Status         string     `gorm:"default:pending;size:20;index" json:"status"`      // pending, paid, shipped, completed, canceled
	PaymentStatus  string     `gorm:"default:unpaid;size:20" json:"payment_status"`     // unpaid, paid, refunded
	ShippingStatus string     `gorm:"default:unshipped;size:20" json:"shipping_status"` // unshipped, shipped, received
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`

	Lines []*OrderLine `json:"lines" gorm:"foreignKey:OrderID"`
}

// OrderLines 订单项
type OrderLine struct {
	ID        string   `gorm:"primaryKey" json:"id"`
	OrderID   string   `gorm:"type:char(36);not null;index" json:"order_id"`
	ProductID string   `gorm:"type:char(36);not null;index" json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductID;references:ID"`
	Quantity  int      `gorm:"not null" json:"quantity"`
	Price     float64  `gorm:"not null" json:"price"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) error {
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	order.Status = consts.ORDER_STATUS_PENDING
	order.PaymentStatus = consts.PAYMENT_STATUS_UNPAID
	order.ShippingStatus = consts.SHIPPING_STATUS_UNSHIPPED
	return nil
}

func (l *OrderLine) BeforeCreate(tx *gorm.DB) error {
	l.ID = uuid.New().String()
	return nil
}

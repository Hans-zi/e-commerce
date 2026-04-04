package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment 支付模型
type Payment struct {
	ID            string     `gorm:"primaryKey" json:"id"`
	OrderID       string     `gorm:"type:char(36);uniqueIndex;not null" json:"order_id"`
	Amount        float64    `gorm:"not null" json:"amount"`
	Method        string     `gorm:"not null;size:20" json:"method"`        // alipay, wechat
	Status        string     `gorm:"default:pending;size:20" json:"status"` // pending, success, failed
	Url           string     `gorm:"not null" json:"url"`
	TransactionID string     `gorm:"size:100" json:"transaction_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().String()

	return nil
}

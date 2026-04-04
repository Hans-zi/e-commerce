package model

import (
	"e-commerce/internal/consts"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Product struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"index;not null;size:200" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null;size:200;" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null;index" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	Status      string         `gorm:"default:active;size:20;index" json:"status"`
	SalesCount  int            `gorm:"default:0" json:"sales_count"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Images      []string       `gorm:"serializer:json" json:"images"`

	CategoryID *string   `gorm:"type:char(36);index" json:"category_id"`
	Category   *Category `json:"category,omitempty"`
}

func (m *Product) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New().String()
	if m.Slug == "" {
		m.Slug = slug.Make(m.Name)
	}
	m.Status = consts.PRODUCT_STATUS_ACTIVE
	return nil
}

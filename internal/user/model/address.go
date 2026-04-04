package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID        string    `gorm:"primaryKey;index" json:"id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"not null;size:50" json:"name"`
	Phone     string    `gorm:"not null;size:50" json:"phone"`
	Province  string    `gorm:"not null;size:50" json:"province"`
	City      string    `gorm:"not null;size:50" json:"city"`
	District  string    `gorm:"not null;size:50" json:"district"`
	Detail    string    `gorm:"not null;size:50" json:"detail"`
	IsDefault bool      `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New().String()
	return nil
}

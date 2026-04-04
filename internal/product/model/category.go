package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Category struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null;size:50"`
	Slug        string `gorm:"uniqueIndex;not null;size:50"`
	Description string `gorm:"not null;size:200"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New().String()
	if c.Slug == "" {
		c.Slug = slug.Make(c.Name)
	}
	return nil
}

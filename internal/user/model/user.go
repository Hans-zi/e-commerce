package model

import (
	"e-commerce/internal/consts"
	"e-commerce/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string          `gorm:"unique;not null;index;primary_key" json:"id"`
	Username  string          `gorm:"not null;size:20;index;unique" json:"username"`
	Password  string          `gorm:"not null" json:"password"`
	Email     string          `gorm:"not null" json:"email"`
	Role      consts.UserRole `gorm:"default:customer" json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `json:"-"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New().String()
	var err error
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	if user.Role == "" {
		user.Role = consts.ROLE_CUSTOMER
	}
	return nil
}

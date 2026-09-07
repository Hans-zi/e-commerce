package repository

import (
	"e-commerce/internal/order/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetUserOrderByID(id, userID string) (*model.Order, error)
	GetByID(id string) (*model.Order, error)
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepo{
		db: db,
	}
}

func (r *orderRepo) GetByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Lines.Product").
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
func (r *orderRepo) GetUserOrderByID(id, userID string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Lines.Product").
		Where("id = ? AND user_id = ?", id, userID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

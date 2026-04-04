package repository

import (
	"e-commerce/internal/order/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *model.Order, lines []*model.OrderLine) error
	GetUserOrderByID(id, userID string) (*model.Order, error)
	GetByID(id string) (*model.Order, error)
	Update(order *model.Order) error
	Delete(id, userID string) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepo{
		db: db,
	}
}

func (r *orderRepo) CreateOrder(order *model.Order, lines []*model.OrderLine) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		for _, line := range lines {
			line.OrderID = order.ID
			err := tx.Create(&line).Error
			if err != nil {
				return err
			}
		}
		return nil
	},
	)

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

func (r *orderRepo) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepo) Delete(id, userID string) error {
	order, err := r.GetUserOrderByID(id, userID)
	if err != nil {
		return err
	}
	return r.db.Select("Lines").Delete(&order).Error
}

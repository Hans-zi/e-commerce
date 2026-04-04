package repository

import (
	"e-commerce/internal/order/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	GetByOrderID(id string) (*model.Payment, error)
	Update(payment *model.Payment) error
}

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepo{
		db: db,
	}
}

func (r *paymentRepo) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepo) GetByOrderID(id string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("order_id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
func (r *paymentRepo) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

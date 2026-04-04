package repository

import (
	"e-commerce/internal/user/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *model.Address) error
	GetByID(id, userID string) (*model.Address, error)
	ListByUser(userID string) ([]*model.Address, error)
	Update(address *model.Address) error
	Delete(id, userID string) error
	SetDefault(id, userID string) error
}

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) AddressRepository {
	return &addressRepo{
		db: db,
	}
}

func (r *addressRepo) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepo) GetByID(id, userID string) (*model.Address, error) {
	var address model.Address
	err := r.db.Where("id = ? and user_id = ?", id, userID).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepo) ListByUser(userID string) ([]*model.Address, error) {
	var addresses []*model.Address
	err := r.db.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepo) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepo) Delete(id, userID string) error {
	if _, err := r.GetByID(id, userID); err != nil {
		return err
	}
	return r.db.Delete(&model.Address{}).Where("id = ? and user_id = ?", id, userID).Error
}

func (r *addressRepo) SetDefault(id, userID string) error {
	address, err := r.GetByID(id, userID)
	if err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&model.Address{}).
			Where("is_default = ?", true).
			Update("is_default", false).Error
		if err != nil {
			return err
		}
		address.IsDefault = true
		err = tx.Save(address).Error
		if err != nil {
			return err
		}
		return nil
	})
}

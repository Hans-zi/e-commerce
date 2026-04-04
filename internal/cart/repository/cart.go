package repository

import (
	"e-commerce/internal/cart/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *model.Cart) error
	GetCartByUserID(userID string) (*model.Cart, error)
	Update(cart *model.Cart) error
	Delete(id string) error
	RemoveCartLine(cartID, productIO string) error
}

type cartRepo struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepo{
		db: db,
	}
}

func (r *cartRepo) Create(cart *model.Cart) error {
	return r.db.Create(&cart).Error
}

func (r *cartRepo) GetCartByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Preload("Lines.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepo) Update(cart *model.Cart) error {
	for _, line := range cart.Lines {
		line.CartID = cart.ID
		product := line.Product
		line.Product = nil
		if err := r.db.Save(line).Error; err != nil {
			return err
		}
		line.Product = product
	}
	return nil
}

func (r *cartRepo) Delete(id string) error {
	return r.db.Delete(&model.Cart{}, id).Error
}

func (r *cartRepo) RemoveCartLine(cartID, productIO string) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productIO).Delete(&model.CartLine{}).Error
}

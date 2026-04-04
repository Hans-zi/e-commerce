package repository

import (
	"e-commerce/internal/consts"
	"e-commerce/internal/order/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductById(id string) (*model.Product, error)
	DecreaseStock(id string, qty int) error
	DecreaseStockBatch(lines []*model.OrderLine) error
	RestoreStockBatch(lines []*model.OrderLine) error
	ListSeckillProduct() ([]*model.Product, error)
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) GetProductById(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) DecreaseStock(id string, qty int) error {
	result := r.db.Where("id = ? AND stock >= ? AND status = ?", id, qty, consts.PRODUCT_STATUS_ACTIVE).
		Model(&model.Product{}).
		UpdateColumn("stock", gorm.Expr("stock - ?", qty))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

func (r *productRepo) RestoreStock(id string, qty int) error {
	result := r.db.Where("id = ?", id).
		Model(&model.Product{}).
		UpdateColumn("stock", gorm.Expr("stock + ?", qty))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

func (r *productRepo) DecreaseStockBatch(lines []*model.OrderLine) error {
	handler := func(tx *gorm.DB) error {
		for _, line := range lines {
			err := r.decreaseStockWithTx(tx, line.ProductID, line.Quantity)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return r.db.Transaction(handler)
}

func (r *productRepo) decreaseStockWithTx(tx *gorm.DB, id string, qty int) error {
	result := tx.Where("id = ? AND stock >= ?", id, qty).
		Model(&model.Product{}).
		UpdateColumn("stock", gorm.Expr("stock - ?", qty))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

func (r *productRepo) RestoreStockBatch(lines []*model.OrderLine) error {
	handler := func(tx *gorm.DB) error {
		for _, line := range lines {
			err := r.restoreStockWithTx(tx, line.ProductID, line.Quantity)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return r.db.Transaction(handler)
}
func (r *productRepo) restoreStockWithTx(tx *gorm.DB, id string, qty int) error {
	result := tx.Where("id = ?", id).
		Model(&model.Product{}).
		UpdateColumn("stock", gorm.Expr("stock + ?", qty))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

func (r *productRepo) ListSeckillProduct() ([]*model.Product, error) {
	db := r.db.Model(&model.Product{})

	orderQuery := fmt.Sprintf("price desc")
	db = db.Order(orderQuery)

	db.Limit(5)

	var list []*model.Product
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

package repository

import (
	"e-commerce/internal/agent/dto"
	"e-commerce/internal/order/model"

	"e-commerce/pkg/utils"
	"fmt"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductById(id string) (*model.Product, error)
	List(req *dto.ListProductsReq) ([]*model.Product, *utils.Pagination, error)
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

func (r *productRepo) List(req *dto.ListProductsReq) ([]*model.Product, *utils.Pagination, error) {
	db := r.db.Model(&model.Product{})
	if req.CategoryID != "" {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Slug != "" {
		db = db.Where("slug LIKE ?", "%"+req.Slug+"%")
	}

	if req.MaxPrice != 0 {
		db = db.Where("price <= ?", req.MaxPrice)
	}
	if req.MinPrice != 0 {
		db = db.Where("price >= ?", req.MinPrice)
	}
	if req.OrderDesc {
		orderQuery := fmt.Sprintf("%s desc", req.OrderBy)
		db = db.Order(orderQuery)
	} else {
		db = db.Order(req.OrderBy)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	paging := utils.NewPaging(req.Page, req.PageSize, total)

	db.Offset(int(paging.Skip)).Limit(int(paging.Limit))

	var list []*model.Product
	if err := db.Find(&list).Error; err != nil {
		return nil, nil, err
	}

	return list, paging, nil
}

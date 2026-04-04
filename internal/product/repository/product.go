package repository

import (
	"e-commerce/internal/product/dto"
	"e-commerce/internal/product/model"
	"e-commerce/pkg/utils"
	"fmt"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(repository *model.Product) error
	GetById(id string) (*model.Product, error)
	GetBySlug(slug string) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id string) error
	List(req *dto.ListProductsReq) ([]*model.Product, *utils.Pagination, error)
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepo) GetById(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *productRepo) GetBySlug(slug string) (*model.Product, error) {
	var product *model.Product
	err := r.db.Where("slug = ?", slug).First(&product).Error
	return product, err
}

func (r *productRepo) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepo) Delete(id string) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepo) List(req *dto.ListProductsReq) ([]*model.Product, *utils.Pagination, error) {
	db := r.db.Model(&model.Product{})
	if req.CategoryID != "" {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", req.Name)
	}
	if req.Slug != "" {
		db = db.Where("slug LIKE ?", req.Slug)
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

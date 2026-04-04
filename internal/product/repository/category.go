package repository

import (
	"e-commerce/internal/product/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	GetById(id string) (*model.Category, error)
	List() ([]*model.Category, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepo) GetById(id string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	return &category, err
}

func (r *categoryRepo) List() ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

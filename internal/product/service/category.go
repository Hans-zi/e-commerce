package service

import (
	"e-commerce/internal/product/dto"
	"e-commerce/internal/product/model"
	"e-commerce/internal/product/repository"
	"e-commerce/pkg/utils"
)

type CategoryService interface {
	Create(req *dto.CreateCategoryReq) (*model.Category, error)
	ListCategories() ([]*model.Category, error)
	GetCategoryByID(id string) (*model.Category, error)
}

type categorySvc struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categorySvc{
		categoryRepo: categoryRepo,
	}
}

func (s *categorySvc) Create(req *dto.CreateCategoryReq) (*model.Category, error) {
	var category model.Category
	utils.Copy(&category, req)
	if err := s.categoryRepo.Create(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *categorySvc) ListCategories() ([]*model.Category, error) {
	return s.categoryRepo.List()
}

func (s *categorySvc) GetCategoryByID(id string) (*model.Category, error) {
	return s.categoryRepo.GetById(id)
}

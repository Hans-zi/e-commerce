package service

import (
	"e-commerce/internal/product/dto"
	"e-commerce/internal/product/model"
	"e-commerce/internal/product/repository"
	"e-commerce/pkg/utils"

	"github.com/bytedance/gopkg/util/logger"
)

type ProductService interface {
	Create(req *dto.CreateProductReq) (*model.Product, error)
	GetProductByID(id string) (*model.Product, error)
	GetProductBySlug(slug string) (*model.Product, error)
	Update(id string, req *dto.UpdateProductReq) (*model.Product, error)
	Delete(id string) error
	ListProducts(req *dto.ListProductsReq) ([]*model.Product, *utils.Pagination, error)
}

type productSvc struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productSvc{productRepo: productRepo}
}

func (s *productSvc) Create(req *dto.CreateProductReq) (*model.Product, error) {
	var product model.Product
	utils.Copy(&product, req)
	err := s.productRepo.Create(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *productSvc) GetProductByID(id string) (*model.Product, error) {
	return s.productRepo.GetById(id)
}

func (s *productSvc) GetProductBySlug(slug string) (*model.Product, error) {
	return s.productRepo.GetBySlug(slug)
}

func (s *productSvc) Update(id string, req *dto.UpdateProductReq) (*model.Product, error) {

	product, err := s.productRepo.GetById(id)
	if err != nil {
		logger.Errorf("Update.GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	logger.Infof("req:  %v", req)
	if err = utils.Copy(product, req); err != nil {
		return nil, err
	}

	err = s.productRepo.Update(product)
	if err != nil {
		logger.Errorf("Update fail, id: %s, error: %s", id, err)
		return nil, err
	}
	return product, err
}

func (s *productSvc) Delete(id string) error {
	return s.productRepo.Delete(id)
}

func (s *productSvc) ListProducts(req *dto.ListProductsReq) ([]*model.Product, *utils.Pagination, error) {
	return s.productRepo.List(req)
}

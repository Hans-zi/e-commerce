package service

import (
	"e-commerce/internal/cart/dto"
	"e-commerce/internal/cart/model"
	"e-commerce/internal/cart/repository"

	"github.com/bytedance/gopkg/util/logger"
)

type CartService interface {
	GetCartByUserID(userID string) (*model.Cart, error)
	AddProduct(userID string, req *dto.AddProductReq) (*model.Cart, error)
	RemoveProduct(userID string, req *dto.RemoveProductReq) (*model.Cart, error)
}

type cartService struct {
	repo repository.CartRepository
}

func NewCartService(cartRepo repository.CartRepository) CartService {
	return &cartService{
		repo: cartRepo,
	}
}

func (s *cartService) GetCartByUserID(userID string) (*model.Cart, error) {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		cart = &model.Cart{
			UserID: userID,
		}
		err = s.repo.Create(cart)
		if err != nil {
			return nil, err
		}
		return cart, err
	}

	return cart, nil
}

func (s *cartService) AddProduct(userID string, req *dto.AddProductReq) (*model.Cart, error) {
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		cart = &model.Cart{
			UserID: userID,
			Lines: []*model.CartLine{{
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}},
		}
		err = s.repo.Create(cart)
		if err != nil {
			return nil, err
		}
	}

	for _, line := range cart.Lines {
		if line.ProductID == req.ProductID {
			line.Quantity = req.Quantity
			err = s.repo.Update(cart)
			if err != nil {
				logger.Errorf("AddProduct.UpdateQuantity fail : %s", err)
				return nil, err
			}
			return cart, nil
		}
	}

	cart.Lines = append(cart.Lines, &model.CartLine{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	err = s.repo.Update(cart)
	if err != nil {
		logger.Errorf("AddProduct.UpdateQuntity fail : %s", err)
		return nil, err
	}
	return s.GetCartByUserID(userID)
}

func (s *cartService) RemoveProduct(userID string, req *dto.RemoveProductReq) (*model.Cart, error) {
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		logger.Errorf("RemoveProduct.GetCartByUserID fail : %s", err)
		return nil, err
	}

	err = s.repo.RemoveCartLine(cart.ID, req.ProductID)
	if err != nil {
		logger.Errorf("RemoveProduct.RemoveCartLine fail : %s", err)
		return nil, err
	}

	return s.GetCartByUserID(userID)
}

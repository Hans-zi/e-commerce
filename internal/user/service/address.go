package service

import (
	"e-commerce/internal/user/dto"
	"e-commerce/internal/user/model"
	"e-commerce/internal/user/repository"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"

	"github.com/bytedance/gopkg/util/logger"
)

type AddressService interface {
	CreateAddress(userID string, req *dto.CreateAddressReq) (*model.Address, error)
	ListAddresses(userID string) ([]*model.Address, error)
	GetAddressByID(id, userID string) (*model.Address, error)
	UpdateAddressByID(id, userID string, req *dto.UpdateAddressReq) (*model.Address, error)
	SetAddressDefault(id, userID string) error
}

type addressSvc struct {
	addressRepo repository.AddressRepository
	maker       token.Maker
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressSvc{
		addressRepo: repo,
	}
}

func (s *addressSvc) CreateAddress(userID string, req *dto.CreateAddressReq) (*model.Address, error) {
	var address model.Address
	utils.Copy(&address, req)
	address.UserID = userID
	err := s.addressRepo.Create(&address)
	if err != nil {
		logger.Fatalf("CreateAddress.Create fail : %s", err)
		return nil, err
	}

	return &address, nil
}

func (s *addressSvc) ListAddresses(userID string) ([]*model.Address, error) {
	return s.addressRepo.ListByUser(userID)
}

func (s *addressSvc) GetAddressByID(id, userID string) (*model.Address, error) {
	return s.addressRepo.GetByID(id, userID)
}

func (s *addressSvc) UpdateAddressByID(id, userID string, req *dto.UpdateAddressReq) (*model.Address, error) {
	address, err := s.addressRepo.GetByID(id, userID)
	if err != nil {
		logger.Fatalf("UpdateAddress.GetCartByUserID fail : %s", err)
		return nil, err
	}
	utils.Copy(address, req)
	err = s.addressRepo.Update(address)
	if err != nil {
		logger.Fatalf("UpdateAddress.Update fail : %s", err)
		return nil, err
	}
	return address, nil
}

func (s *addressSvc) SetAddressDefault(id, userID string) error {
	return s.addressRepo.SetDefault(id, userID)
}

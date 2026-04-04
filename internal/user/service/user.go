package service

import (
	"e-commerce/internal/user/dto"
	"e-commerce/internal/user/model"
	"e-commerce/internal/user/repository"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"

	"github.com/bytedance/gopkg/util/logger"
)

type UserService interface {
	Register(req *dto.RegisterReq) (*model.User, error)
	Login(req *dto.LoginReq) (*model.User, string, error)
	GetUserById(id string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	ChangePassword(id string, req *dto.ChangePasswordReq) (*model.User, error)
}

type userSvc struct {
	userRepo repository.UserRepository
	maker    token.Maker
}

func NewUserService(repo repository.UserRepository, tokenMaker token.Maker) UserService {
	return &userSvc{
		userRepo: repo,
		maker:    tokenMaker,
	}
}

func (s *userSvc) Register(req *dto.RegisterReq) (*model.User, error) {
	var user model.User
	utils.Copy(&user, &req)
	err := s.userRepo.Create(&user)
	if err != nil {
		logger.Errorf("Register.Fail, error: %s", err)
		return nil, err
	}
	return &user, nil
}

func (s *userSvc) Login(req *dto.LoginReq) (*model.User, string, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		logger.Errorf("Login.GetByUsername fail: %s", err)
		return nil, "", err
	}

	err = utils.CheckHashPassWord(user.Password, req.Password)
	if err != nil {
		logger.Errorf("Login.CheckPassword fail: %s", err)
		return nil, "", err
	}

	jtoken, err := s.maker.CreateToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Errorf("Login.CreateToken fail: %s", err)
		return nil, "", err
	}
	return user, jtoken, nil
}
func (s *userSvc) GetUserById(id string) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userSvc) GetByUsername(username string) (*model.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *userSvc) ChangePassword(id string, req *dto.ChangePasswordReq) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Fatalf("ChangePassword.GetUserById fail : %s", err)
		return nil, err
	}
	err = utils.CheckHashPassWord(user.Password, req.Password)
	if err != nil {
		logger.Fatalf("ChangePassword.CheckPassword fail : %s", err)
		return nil, err
	}
	newHashPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logger.Fatalf("ChangePassword.HashPassword fail : %s", err)
		return nil, err
	}
	user.Password = newHashPassword
	err = s.userRepo.Update(user)
	if err != nil {
		logger.Fatalf("ChangePassword.Update ail : %s", err)
		return nil, err
	}
	return user, nil
}

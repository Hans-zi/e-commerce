package repository

import (
	"e-commerce/internal/user/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	GetByID(id string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}
func (r *userRepo) GetByID(id string) (*model.User, error) {
	args := model.User{
		ID: id,
	}
	var user model.User
	res := r.db.Find(&args).First(&user)
	return &user, res.Error
}

func (r *userRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	res := r.db.Where("username = ?", username).First(&user)
	return &user, res.Error
}

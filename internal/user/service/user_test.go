package service

import (
	"e-commerce/internal/user/model"
	"e-commerce/internal/user/repository"
	"e-commerce/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) model.User {
	//password := utils.RandomString(6)
	password := "xiaohan1234"
	user := model.User{
		Username:  utils.RandomString(6),
		Password:  password,
		Email:     utils.RandomEmail(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userService := NewUserService(repository.NewUserRepo(testDB))
	err := userService.Create(&user)
	require.NoError(t, err)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserById(t *testing.T) {
	userService := NewUserService(repository.NewUserRepo(testDB))
	user1 := createRandomUser(t)
	user2, err := userService.GetById(user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.Email, user1.Email)
}

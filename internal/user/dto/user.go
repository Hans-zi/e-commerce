package dto

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=5"`
	Email    string `json:"email" binding:"required,email"`
}

type RegisterRes struct {
	User User `json:"user"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=5"`
}

type LoginRes struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ChangePasswordReq struct {
	Password    string `json:"password" binding:"required,min=5"`
	NewPassword string `json:"new_password" binding:"required,min=5"`
}

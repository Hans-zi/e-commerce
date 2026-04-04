package token

import "e-commerce/internal/consts"

type Maker interface {
	CreateToken(userID string, username string, role consts.UserRole) (string, error)
	VerifyToken(token string) (*Payload, error)
}

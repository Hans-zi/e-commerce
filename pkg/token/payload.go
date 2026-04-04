package token

import (
	"e-commerce/internal/consts"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID   string
	Username string
	Role     consts.UserRole

	jwt.RegisteredClaims
}

func NewPayload(userID string, username string, role consts.UserRole, duration time.Duration) *Payload {
	return &Payload{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
}

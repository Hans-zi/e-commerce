package token

import (
	"e-commerce/internal/consts"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	SymmetricKey string
	Duration     time.Duration
}

const symmetricKeySize = 32

var (
	ErrExpiredToken = errors.New("token is Expired")
	ErrInvalidToken = errors.New("token is Invalid")
)

func NewJWTMaker(symmetricKey string, duration time.Duration) (Maker, error) {
	if len(symmetricKey) != symmetricKeySize {
		err := errors.New(fmt.Sprintf("Too short key, want %d characters", symmetricKeySize))
		return nil, err
	}
	return &JWTMaker{
		SymmetricKey: symmetricKey,
		Duration:     duration,
	}, nil
}

func (maker *JWTMaker) CreateToken(userID string, username string, role consts.UserRole) (string, error) {
	payload := NewPayload(userID, username, role, maker.Duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(maker.SymmetricKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (any, error) {
		return []byte(maker.SymmetricKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

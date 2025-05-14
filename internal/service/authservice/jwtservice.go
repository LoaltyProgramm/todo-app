package authservice

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	CreateToken(password string) (string, error)
	ParseToken(token string) (bool, error)
}

var secretKey = []byte("sdfslkdnf2893jkfns893")

type JWTServiceImpl struct{}

func NewJWTService() JWTService {
	return &JWTServiceImpl{}
}

func (s *JWTServiceImpl) CreateToken(password string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("token is not signed: %w", err)
	}

	return signedToken, nil
}

func (s *JWTServiceImpl) ParseToken(tokenStr string) (bool, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid{
		return false, fmt.Errorf("parse token fail: %w", err)
	}

	return token.Valid, nil		
}

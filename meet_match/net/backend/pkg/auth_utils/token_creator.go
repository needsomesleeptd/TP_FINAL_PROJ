package auth_utils

import (
	"errors"
	"fmt"
	"test_backend_frontend/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ITokenHandler interface {
	GenerateToken(credentials models.User, key string) (string, error)
	ValidateToken(tokenString string, key string) error
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrParsingToken = errors.New("error parsing token")
)

type JWTTokenHandler struct {
}

func NewJWTTokenHandler() ITokenHandler {
	return JWTTokenHandler{}
}

func (hasher JWTTokenHandler) GenerateToken(credentials models.User, key string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exprires": time.Now().Add(time.Hour * 24),
			"login":    credentials.Login,
			"ID":       credentials.ID,
		})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("creating token err: %w", err)
	}

	return tokenString, nil
}

func (hasher JWTTokenHandler) ValidateToken(tokenString string, key string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return ErrParsingToken
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

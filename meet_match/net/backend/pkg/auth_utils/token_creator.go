package auth_utils

import (
	"errors"
	"fmt"
	"test_backend_frontend/internal/models"
	"time"

	errors2 "github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
)

type Payload struct {
	Login string
	ID    uint64
	Exp   time.Duration
}

type ITokenHandler interface {
	GenerateToken(credentials models.User, key string) (string, error)
	ValidateToken(tokenString string, key string) error
	ParseToken(tokenString string, key string) (*Payload, error)
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
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = credentials.Login
	claims["id"] = credentials.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("creating token err: %w", err)
	}

	return tokenString, nil
}

func (hasher JWTTokenHandler) ValidateToken(tokenString string, key string) error {
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return ErrParsingToken
	}

	if !jwtToken.Valid {
		return ErrInvalidToken
	}

	return nil
}

func (hasher JWTTokenHandler) ParseToken(tokenString string, key string) (*Payload, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, errors2.Wrap(err, "auth.tokenhelper.ParseToken error in parse")
	}

	payload := &Payload{
		Login: claims["login"].(string),
		ID:    uint64(claims["id"].(float64)),
		Exp:   time.Duration(claims["exp"].(float64)),
	}

	return payload, nil
}

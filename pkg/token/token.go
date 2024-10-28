package token

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/config"
	"apiGateway/pkg/logger"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// tokenClaims - структура токена JWT
type tokenClaims struct {
	jwt.StandardClaims
	UserId uint64 `json:"userId"`
	Role   string `json:"role"`
}

// IUser - интерфейс для доступа к полям субъекта токена JWT
type IUser interface {
	GetUserId() uint64
	GetRole() string
}

// GetUserId - возвращает ID пользователя из токена JWT
func (t *tokenClaims) GetUserId() uint64 {
	return t.UserId
}

// GetRole - возвращает роль юзера из токена JWT
func (t *tokenClaims) GetRole() string {
	return t.Role
}

// CreateToken - создание токена JWT
func CreateToken(user *DatabaseServicev1.CreateUserResponse, cfg *config.Config) (string, error) {
	parsedValue, err := time.ParseDuration(cfg.Jwt.Expires)
	if err != nil {
		logger.Error("Ошибка при парсинге времени жизни токена: %v", err)
		return "", err
	}

	standartClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(parsedValue).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	claims := &tokenClaims{
		StandardClaims: standartClaims,
		UserId:         user.GetId(),
		Role:           user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.Jwt.Secret))
	if err != nil {
		logger.Error("Ошибка при подписи токена: %v", err)
		return "", err
	}

	return signedToken, nil
}

// ParseToken - парсит токен из строки
func ParseToken(accessToken string, cfg *config.Config) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprint("Неверная подпись"))
		}
		return []byte(cfg.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}

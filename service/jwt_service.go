package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"gitlab.com/katsuotz/skip-api/entity"
	"os"
)

type JWTService interface {
	GenerateToken(ctx context.Context, request entity.User) string
}

type jwtService struct {
	secretKey string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "PasarUdang V.02"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(ctx context.Context, request entity.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": request.UserID,
		"role_id": request.Role,
	})
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

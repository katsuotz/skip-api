package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gitlab.com/katsuotz/skip-api/entity"
	"net/http"
	"os"
	"strings"
)

type JWTService interface {
	GenerateToken(ctx context.Context, request entity.User) string
	ValidateToken(token string) (*jwt.Token, error)
	IsLoggedIn(ctx *gin.Context)
	IsGuest(ctx *gin.Context)
	IsAdmin(ctx *gin.Context)
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
		"user_id": request.ID,
		"role":    request.Role,
	})
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected login method %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}

func (s *jwtService) IsLoggedIn(ctx *gin.Context) {
	authorization := ctx.Request.Header["Authorization"]

	if authorization != nil {
		token := strings.Split(authorization[0], " ")[1]
		aToken, err := s.ValidateToken(token)
		if err == nil {
			claims := aToken.Claims.(jwt.MapClaims)
			if claims["user_id"] != nil {
				ctx.Set("user_id", claims["user_id"])
				ctx.Set("role", claims["role"])
				ctx.Next()
				return
			}
		}
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

func (s *jwtService) IsGuest(ctx *gin.Context) {
	authorization := ctx.Request.Header["Authorization"]

	if authorization == nil {
		ctx.Next()
		return
	}

	token := strings.Split(authorization[0], " ")[1]
	aToken, err := s.ValidateToken(token)
	if err == nil {
		claims := aToken.Claims.(jwt.MapClaims)
		if claims["user_id"] == nil {
			ctx.Next()
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

func (s *jwtService) IsAdmin(ctx *gin.Context) {
	role := ctx.MustGet("role")
	if role == "admin" {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

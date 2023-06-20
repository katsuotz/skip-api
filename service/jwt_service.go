package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/helper"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

type JWTService interface {
	GenerateToken(ctx context.Context, user dto.UserResponse) string
	ValidateToken(token string) (*jwt.Token, error)
	IsLoggedIn(ctx *gin.Context)
	IsGuest(ctx *gin.Context)
	IsAdmin(ctx *gin.Context)
	IsPegawai(ctx *gin.Context)
	IsStaff(ctx *gin.Context)
	IsNotSiswa(ctx *gin.Context)
	IsRole(ctx *gin.Context, roles string)
}

type jwtService struct {
	secretKey string
	db        *gorm.DB
}

func NewJWTService(db *gorm.DB) JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		db:        db,
	}
}

func getSecretKey() string {
	return os.Getenv("JWT_SECRET")
}

func (s *jwtService) GenerateToken(ctx context.Context, user dto.UserResponse) string {
	jwtData := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
	}

	if user.PegawaiID != 0 {
		jwtData["pegawai_id"] = user.PegawaiID
	} else if user.Role == "siswa" {
		jwtData["siswa_id"] = user.SiswaID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	tokenString, err := token.SignedString([]byte(s.secretKey))
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
				ctx.Set("pegawai_id", claims["pegawai_id"])
				ctx.Set("siswa_id", claims["siswa_id"])
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
	if role == "admin" || role == "staff-ict" {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

func (s *jwtService) IsPegawai(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)
	if helper.IsPegawai(role) {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

func (s *jwtService) IsStaff(ctx *gin.Context) {
	roles := "admin;staff-ict;guru-bk;"
	role := ctx.MustGet("role").(string) + ";"
	if strings.Contains(roles, role) {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

func (s *jwtService) IsNotSiswa(ctx *gin.Context) {
	role := ctx.MustGet("role")
	if role != "siswa" {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

func (s *jwtService) IsRole(ctx *gin.Context, roles string) {
	role := ctx.MustGet("role").(string)
	if strings.Contains(roles, role) {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

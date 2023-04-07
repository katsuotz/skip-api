package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

type JWTService interface {
	GenerateToken(ctx context.Context, user entity.User) string
	ValidateToken(token string) (*jwt.Token, error)
	IsLoggedIn(ctx *gin.Context)
	IsGuest(ctx *gin.Context)
	IsAdmin(ctx *gin.Context)
	IsGuru(ctx *gin.Context)
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

func (s *jwtService) GenerateToken(ctx context.Context, user entity.User) string {
	jwtData := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
	}

	if user.Role == "guru" {
		var guru entity.Guru
		s.db.Where("user_id = ?", user.ID).First(&guru)
		jwtData["guru_id"] = guru.ID
	}

	if user.Role == "siswa" {
		var siswa entity.Siswa
		s.db.Where("user_id = ?", user.ID).First(&siswa)
		jwtData["siswa_id"] = siswa.ID
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
				ctx.Set("guru_id", claims["guru_id"])
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
	if role == "admin" {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

func (s *jwtService) IsGuru(ctx *gin.Context) {
	roles := []string{"guru", "staff-ict", "guru-bk", "tata-usaha"}

	role := ctx.MustGet("role").(string)
	guruID := int(ctx.MustGet("guru_id").(float64))
	if strings.Contains(strings.Join(roles, ","), role) && guruID != 0 {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

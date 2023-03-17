package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/service"
	"net/http"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Tes(ctx *gin.Context)
}

type authController struct {
	UserRepository repository.UserRepository
	JWTService     service.JWTService
}

func NewAuthController(userRepository repository.UserRepository, jwtService service.JWTService) AuthController {
	return &authController{
		userRepository,
		jwtService,
	}
}

func (c *authController) Tes(ctx *gin.Context) {
	ctx.Data(200, "application/json; charset=utf-8", []byte("connected"))
	return
}
func (c *authController) Login(ctx *gin.Context) {
	loginReq := dto.LoginRequest{}
	errDTO := ctx.ShouldBindJSON(&loginReq)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	user := c.UserRepository.FindByUsername(ctx, loginReq.Username)
	match := helper.CheckPasswordHash(loginReq.Password, user.Password)
	if !match || user.UserID == 0 {
		response := helper.BuildErrorResponse("Wrong username or password", nil, nil)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}
	token := c.JWTService.GenerateToken(ctx, user)
	if token != "" {
		loginRes := dto.LoginResponse{}
		loginRes.User = user
		loginRes.Token = token
		response := helper.BuildSuccessResponse("Login Success", loginRes)
		ctx.JSON(http.StatusOK, response)
	}
}

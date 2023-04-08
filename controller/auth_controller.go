package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/service"
	"net/http"
	"strconv"
)

type AuthController interface {
	Login(ctx *gin.Context)
	GetLog(ctx *gin.Context)
	Tes(ctx *gin.Context)
}

type authController struct {
	UserRepository     repository.UserRepository
	LoginLogRepository repository.LoginLogRepository
	JWTService         service.JWTService
}

func NewAuthController(
	userRepository repository.UserRepository,
	loginLogRepository repository.LoginLogRepository,
	jwtService service.JWTService,
) AuthController {
	return &authController{
		userRepository,
		loginLogRepository,
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
	if !match || user.ID == 0 {
		if user.ID != 0 {
			go c.UserRepository.LoginLog(ctx, user.ID, "Failed Login Attempt")
		}
		response := helper.BuildErrorResponse("Wrong username or password", nil, nil)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}
	token := c.JWTService.GenerateToken(ctx, user)
	go c.UserRepository.LoginLog(ctx, user.ID, "Successful Login")

	loginRes := dto.LoginResponse{}
	loginRes.User = user
	loginRes.Token = token
	response := helper.BuildSuccessResponse("Login Success", loginRes)
	ctx.JSON(http.StatusOK, response)
}

func (c *authController) GetLog(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	search := ctx.DefaultQuery("search", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	guru := c.LoginLogRepository.GetLog(ctx, pageInt, perPageInt, search)
	response := helper.BuildSuccessResponse("", guru)
	ctx.JSON(http.StatusOK, response)
	return
}

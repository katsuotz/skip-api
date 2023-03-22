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

type GuruController interface {
	GetGuru(ctx *gin.Context)
	CreateGuru(ctx *gin.Context)
	UpdateGuru(ctx *gin.Context)
	DeleteGuru(ctx *gin.Context)
}

type guruController struct {
	GuruRepository repository.GuruRepository
	JWTService     service.JWTService
}

func NewGuruController(guruRepository repository.GuruRepository, jwtService service.JWTService) GuruController {
	return &guruController{
		guruRepository,
		jwtService,
	}
}

func (c *guruController) GetGuru(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	search := ctx.DefaultQuery("search", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	guru := c.GuruRepository.GetGuru(ctx, pageInt, perPageInt, search)
	response := helper.BuildSuccessResponse("", guru)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *guruController) CreateGuru(ctx *gin.Context) {
	req := dto.GuruRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.GuruRepository.CreateGuru(ctx, req)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Guru created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *guruController) UpdateGuru(ctx *gin.Context) {
	req := dto.GuruRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	guruID, err := strconv.Atoi(ctx.Param("guru_id"))
	if err != nil || guruID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.GuruRepository.UpdateGuru(ctx, req, guruID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Guru created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *guruController) DeleteGuru(ctx *gin.Context) {
	guruID, err := strconv.Atoi(ctx.Param("guru_id"))
	if err != nil || guruID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.GuruRepository.DeleteGuru(ctx, int(guruID))

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Guru deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

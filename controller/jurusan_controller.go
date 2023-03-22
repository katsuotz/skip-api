package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/service"
	"net/http"
	"strconv"
)

type JurusanController interface {
	GetJurusan(ctx *gin.Context)
	CreateJurusan(ctx *gin.Context)
	UpdateJurusan(ctx *gin.Context)
	DeleteJurusan(ctx *gin.Context)
}

type jurusanController struct {
	JurusanRepository repository.JurusanRepository
	JWTService        service.JWTService
}

func NewJurusanController(jurusanRepository repository.JurusanRepository, jwtService service.JWTService) JurusanController {
	return &jurusanController{
		jurusanRepository,
		jwtService,
	}
}

func (c *jurusanController) GetJurusan(ctx *gin.Context) {
	jurusan := c.JurusanRepository.GetJurusan(ctx)
	response := helper.BuildSuccessResponse("", jurusan)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *jurusanController) CreateJurusan(ctx *gin.Context) {
	req := dto.JurusanRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	jurusan := entity.Jurusan{
		NamaJurusan: req.NamaJurusan,
	}

	_, err := c.JurusanRepository.CreateJurusan(ctx, jurusan)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Jurusan created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *jurusanController) UpdateJurusan(ctx *gin.Context) {
	req := dto.JurusanRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	jurusanID, err := strconv.Atoi(ctx.Param("jurusan_id"))
	if err != nil || jurusanID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newJurusan := entity.Jurusan{
		ID:          jurusanID,
		NamaJurusan: req.NamaJurusan,
	}

	_, err = c.JurusanRepository.UpdateJurusan(ctx, newJurusan)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Jurusan updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *jurusanController) DeleteJurusan(ctx *gin.Context) {
	jurusanID, err := strconv.Atoi(ctx.Param("jurusan_id"))
	if err != nil || jurusanID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.JurusanRepository.DeleteJurusan(ctx, jurusanID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Jurusan deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

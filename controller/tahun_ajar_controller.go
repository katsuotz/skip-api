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

type TahunAjarController interface {
	GetTahunAjar(ctx *gin.Context)
	CreateTahunAjar(ctx *gin.Context)
	UpdateTahunAjar(ctx *gin.Context)
	DeleteTahunAjar(ctx *gin.Context)
	SetActiveTahunAjar(ctx *gin.Context)
}

type tahunAjarController struct {
	TahunAjarRepository repository.TahunAjarRepository
	JWTService          service.JWTService
}

func NewTahunAjarController(tahunAjarRepository repository.TahunAjarRepository, jwtService service.JWTService) TahunAjarController {
	return &tahunAjarController{
		tahunAjarRepository,
		jwtService,
	}
}

func (c *tahunAjarController) GetTahunAjar(ctx *gin.Context) {
	tahunAjar := c.TahunAjarRepository.GetTahunAjar(ctx)
	response := helper.BuildSuccessResponse("", tahunAjar)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) CreateTahunAjar(ctx *gin.Context) {
	req := dto.TahunAjarRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	tahunAjar := entity.TahunAjar{
		TahunAjar: req.TahunAjar,
		IsActive:  false,
	}

	_, err := c.TahunAjarRepository.CreateTahunAjar(ctx, tahunAjar)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) UpdateTahunAjar(ctx *gin.Context) {
	req := dto.TahunAjarRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	tahunAjarID, err := strconv.Atoi(ctx.Param("tahun_ajar_id"))
	if err != nil || tahunAjarID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newTahunAjar := entity.TahunAjar{
		ID:        tahunAjarID,
		TahunAjar: req.TahunAjar,
	}

	_, err = c.TahunAjarRepository.UpdateTahunAjar(ctx, newTahunAjar)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) DeleteTahunAjar(ctx *gin.Context) {
	tahunAjarID, err := strconv.Atoi(ctx.Param("tahun_ajar_id"))
	if err != nil || tahunAjarID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.TahunAjarRepository.DeleteTahunAjar(ctx, int(tahunAjarID))

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) SetActiveTahunAjar(ctx *gin.Context) {
	tahunAjarID, err := strconv.Atoi(ctx.Param("tahun_ajar_id"))
	if err != nil || tahunAjarID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.TahunAjarRepository.SetActiveTahunAjar(ctx, int(tahunAjarID))

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

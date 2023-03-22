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

type SiswaController interface {
	GetSiswa(ctx *gin.Context)
	CreateSiswa(ctx *gin.Context)
	UpdateSiswa(ctx *gin.Context)
	DeleteSiswa(ctx *gin.Context)
}

type siswaController struct {
	SiswaRepository repository.SiswaRepository
	JWTService      service.JWTService
}

func NewSiswaController(siswaRepository repository.SiswaRepository, jwtService service.JWTService) SiswaController {
	return &siswaController{
		siswaRepository,
		jwtService,
	}
}

func (c *siswaController) GetSiswa(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	kelasID := ctx.DefaultQuery("kelas_id", "0")
	search := ctx.DefaultQuery("search", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	kelasIDInt, _ := strconv.Atoi(kelasID)

	siswa := c.SiswaRepository.GetSiswa(ctx, pageInt, perPageInt, search, kelasIDInt)
	response := helper.BuildSuccessResponse("", siswa)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) CreateSiswa(ctx *gin.Context) {
	req := dto.SiswaRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.SiswaRepository.CreateSiswa(ctx, req)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) UpdateSiswa(ctx *gin.Context) {
	req := dto.SiswaRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	siswaID, err := strconv.Atoi(ctx.Param("siswa_id"))
	if err != nil || siswaID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.SiswaRepository.UpdateSiswa(ctx, req, siswaID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) DeleteSiswa(ctx *gin.Context) {
	siswaID, err := strconv.Atoi(ctx.Param("siswa_id"))
	if err != nil || siswaID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.SiswaRepository.DeleteSiswa(ctx, int(siswaID))

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
	"strconv"
)

type SiswaController interface {
	GetSiswa(ctx *gin.Context)
	GetSiswaDetailByNIS(ctx *gin.Context)
	CreateSiswa(ctx *gin.Context)
	UpdateSiswa(ctx *gin.Context)
	DeleteSiswa(ctx *gin.Context)
}

type siswaController struct {
	SiswaRepository     repository.SiswaRepository
	PoinSiswaRepository repository.PoinSiswaRepository
}

func NewSiswaController(
	siswaRepository repository.SiswaRepository,
	poinSiswaRepository repository.PoinSiswaRepository,
) SiswaController {
	return &siswaController{
		siswaRepository,
		poinSiswaRepository,
	}
}

func (c *siswaController) GetSiswa(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	kelasID := ctx.DefaultQuery("kelas_id", "")
	search := ctx.DefaultQuery("search", "")
	tahunAjarActive := ctx.DefaultQuery("tahun_ajar_active", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	siswa := c.SiswaRepository.GetSiswa(ctx, pageInt, perPageInt, search, kelasID, tahunAjarActive)
	response := helper.BuildSuccessResponse("", siswa)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) GetSiswaDetailByNIS(ctx *gin.Context) {
	nis := ctx.Param("nis")
	if nis == "" {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result := dto.SiswaDetailLog{}

	result.Siswa = c.SiswaRepository.GetSiswaByNIS(ctx, nis)
	result.Log = c.PoinSiswaRepository.GetPoinLogSiswaByKelas(ctx, nis)

	response := helper.BuildSuccessResponse("", result)
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

	err = c.SiswaRepository.DeleteSiswa(ctx, siswaID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

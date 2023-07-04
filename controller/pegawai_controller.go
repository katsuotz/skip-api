package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
	"strconv"
)

type PegawaiController interface {
	GetPegawai(ctx *gin.Context)
	CreatePegawai(ctx *gin.Context)
	UpdatePegawai(ctx *gin.Context)
	DeletePegawai(ctx *gin.Context)
}

type pegawaiController struct {
	PegawaiRepository repository.PegawaiRepository
}

func NewPegawaiController(pegawaiRepository repository.PegawaiRepository) PegawaiController {
	return &pegawaiController{
		pegawaiRepository,
	}
}

func (c *pegawaiController) GetPegawai(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	search := ctx.DefaultQuery("search", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	pegawai := c.PegawaiRepository.GetPegawai(ctx, pageInt, perPageInt, search)
	response := helper.BuildSuccessResponse("", pegawai)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *pegawaiController) CreatePegawai(ctx *gin.Context) {
	req := dto.CreatePegawaiRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.PegawaiRepository.CreatePegawai(ctx, req)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Pegawai berhasil dibuat", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *pegawaiController) UpdatePegawai(ctx *gin.Context) {
	req := dto.PegawaiRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pegawaiID, err := strconv.Atoi(ctx.Param("pegawai_id"))
	if err != nil || pegawaiID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.PegawaiRepository.UpdatePegawai(ctx, req, pegawaiID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Pegawai berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *pegawaiController) DeletePegawai(ctx *gin.Context) {
	pegawaiID, err := strconv.Atoi(ctx.Param("pegawai_id"))
	if err != nil || pegawaiID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.PegawaiRepository.DeletePegawai(ctx, pegawaiID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Pegawai berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

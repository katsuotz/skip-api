package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
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
}

func NewTahunAjarController(tahunAjarRepository repository.TahunAjarRepository) TahunAjarController {
	return &tahunAjarController{
		tahunAjarRepository,
	}
}

func (c *tahunAjarController) GetTahunAjar(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	search := ctx.DefaultQuery("search", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	tahunAjar := c.TahunAjarRepository.GetTahunAjar(ctx, pageInt, perPageInt, search)
	response := helper.BuildSuccessResponse("", tahunAjar)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) CreateTahunAjar(ctx *gin.Context) {
	req := dto.TahunAjarRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	tahunAjar := entity.TahunAjar{
		TahunAjar: req.TahunAjar,
		IsActive:  false,
	}

	_, err := c.TahunAjarRepository.CreateTahunAjar(ctx, tahunAjar)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar berhasil dibuat", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) UpdateTahunAjar(ctx *gin.Context) {
	req := dto.TahunAjarRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	tahunAjarID, err := strconv.Atoi(ctx.Param("tahun_ajar_id"))
	if err != nil || tahunAjarID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newTahunAjar := entity.TahunAjar{
		ID:        tahunAjarID,
		TahunAjar: req.TahunAjar,
	}

	_, err = c.TahunAjarRepository.UpdateTahunAjar(ctx, newTahunAjar)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) DeleteTahunAjar(ctx *gin.Context) {
	tahunAjarID, err := strconv.Atoi(ctx.Param("tahun_ajar_id"))
	if err != nil || tahunAjarID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.TahunAjarRepository.DeleteTahunAjar(ctx, tahunAjarID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *tahunAjarController) SetActiveTahunAjar(ctx *gin.Context) {
	tahunAjarID, err := strconv.Atoi(ctx.Param("tahun_ajar_id"))
	if err != nil || tahunAjarID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.TahunAjarRepository.SetActiveTahunAjar(ctx, tahunAjarID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Tahun Ajar berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

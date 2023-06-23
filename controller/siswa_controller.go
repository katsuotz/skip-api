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
	SiswaRepository   repository.SiswaRepository
	PoinLogRepository repository.PoinLogRepository
}

func NewSiswaController(
	siswaRepository repository.SiswaRepository,
	poinLogRepository repository.PoinLogRepository,
) SiswaController {
	return &siswaController{
		siswaRepository,
		poinLogRepository,
	}
}

func (c *siswaController) GetSiswa(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	kelasID := ctx.DefaultQuery("kelas_id", "")
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	jurusanID := ctx.DefaultQuery("jurusan_id", "")
	search := ctx.DefaultQuery("search", "")
	tahunAjarActive := ctx.DefaultQuery("tahun_ajar_active", "")
	summary := ctx.DefaultQuery("summary", "")
	summaryBool, _ := strconv.ParseBool(summary)
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	siswa := c.SiswaRepository.GetSiswa(ctx, pageInt, perPageInt, search, kelasID, tahunAjarID, jurusanID, tahunAjarActive, summaryBool)
	response := helper.BuildSuccessResponse("", siswa)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) GetSiswaDetailByNIS(ctx *gin.Context) {
	nis := ctx.Param("nis")
	if nis == "" {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result := dto.SiswaDetailLog{}

	result.Siswa = c.SiswaRepository.GetSiswaByNIS(ctx, nis)

	//role := ctx.MustGet("role")

	//if role == "siswa" {
	//	siswaID := int(ctx.MustGet("siswa_id").(float64))
	//
	//	if siswaID != result.Siswa.ID {
	//		response := helper.BuildErrorResponse("Unauthorized", nil, nil)
	//		ctx.JSON(http.StatusUnauthorized, response)
	//		return
	//	}
	//}

	result.Log = c.PoinLogRepository.GetPoinLogSiswaByKelas(ctx, nis)

	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) CreateSiswa(ctx *gin.Context) {
	req := dto.SiswaRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.SiswaRepository.CreateSiswa(ctx, req)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa berhasil dibuat", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) UpdateSiswa(ctx *gin.Context) {
	req := dto.SiswaRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	siswaID, err := strconv.Atoi(ctx.Param("siswa_id"))
	if err != nil || siswaID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.SiswaRepository.UpdateSiswa(ctx, req, siswaID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa berhasil dibuat", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *siswaController) DeleteSiswa(ctx *gin.Context) {
	siswaID, err := strconv.Atoi(ctx.Param("siswa_id"))
	if err != nil || siswaID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.SiswaRepository.DeleteSiswa(ctx, siswaID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

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

type PoinSiswaController interface {
	GetPoinSiswa(ctx *gin.Context)
	GetPoinKelas(ctx *gin.Context)
	GetPoinJurusan(ctx *gin.Context)
	AddPoinSiswa(ctx *gin.Context)
	UpdatePoinSiswa(ctx *gin.Context)
	DeletePoinSiswa(ctx *gin.Context)
}

type poinSiswaController struct {
	PoinSiswaRepository repository.PoinSiswaRepository
	PoinLogRepository   repository.PoinLogRepository
}

func NewPoinSiswaController(
	poinSiswaRepository repository.PoinSiswaRepository,
	poinLogRepository repository.PoinLogRepository,
) PoinSiswaController {
	return &poinSiswaController{
		poinSiswaRepository,
		poinLogRepository,
	}
}

func (c *poinSiswaController) GetPoinSiswa(ctx *gin.Context) {
	siswaKelasID, err := strconv.Atoi(ctx.Param("siswa_kelas_id"))
	if err != nil || siswaKelasID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	poinSiswa := c.PoinSiswaRepository.GetPoinSiswa(ctx, siswaKelasID)
	response := helper.BuildSuccessResponse("", poinSiswa)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *poinSiswaController) GetPoinKelas(ctx *gin.Context) {
	kelasID, err := strconv.Atoi(ctx.Param("kelas_id"))
	if err != nil || kelasID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	poinSiswa := c.PoinSiswaRepository.GetPoinKelas(ctx, kelasID)
	response := helper.BuildSuccessResponse("", poinSiswa)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *poinSiswaController) GetPoinJurusan(ctx *gin.Context) {
	jurusanID := ctx.Param("jurusan_id")
	if jurusanID == "" {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	tahunAjarID := ctx.Query("tahun_ajar_id")
	if tahunAjarID == "" {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	poinSiswa := c.PoinSiswaRepository.GetPoinJurusan(ctx, jurusanID, tahunAjarID)
	response := helper.BuildSuccessResponse("", poinSiswa)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *poinSiswaController) AddPoinSiswa(ctx *gin.Context) {
	req := dto.PoinSiswaRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pegawaiID := int(ctx.MustGet("pegawai_id").(float64))
	req.PegawaiID = pegawaiID

	err := c.PoinSiswaRepository.AddPoinSiswa(ctx, req)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Poin berhasil dibuat", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *poinSiswaController) UpdatePoinSiswa(ctx *gin.Context) {
	req := dto.UpdatePoinLogRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	poinLogID, err := strconv.Atoi(ctx.Param("poin_log_id"))
	if err != nil || poinLogID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	poinLog := entity.PoinLog{
		ID:           poinLogID,
		TindakLanjut: req.TindakLanjut,
	}

	err = c.PoinSiswaRepository.UpdatePoinSiswa(ctx, poinLog)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Poin berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *poinSiswaController) DeletePoinSiswa(ctx *gin.Context) {
	poinLogID, err := strconv.Atoi(ctx.Param("poin_log_id"))
	if err != nil || poinLogID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.PoinSiswaRepository.DeletePoinSiswa(ctx, poinLogID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Poin berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

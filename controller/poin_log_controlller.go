package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
	"strconv"
)

type PoinLogController interface {
	GetPoinLog(ctx *gin.Context)
	GetPoinSiswaLog(ctx *gin.Context)
}

type poinLogController struct {
	PoinLogRepository repository.PoinLogRepository
}

func NewPoinLogController(
	poinLogRepository repository.PoinLogRepository,
) PoinLogController {
	return &poinLogController{
		poinLogRepository,
	}
}

func (c *poinLogController) GetPoinLog(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	order := ctx.DefaultQuery("order", "asc")
	orderBy := ctx.DefaultQuery("order_by", "nama")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinLogRepository.GetPoinLogPagination(ctx, pageInt, perPageInt, order, orderBy, tahunAjarID, pegawaiID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *poinLogController) GetPoinSiswaLog(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	siswaKelasID, err := strconv.Atoi(ctx.Param("siswa_kelas_id"))
	if err != nil || siswaKelasID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	poinSiswa := c.PoinLogRepository.GetPoinSiswaLog(ctx, pageInt, perPageInt, siswaKelasID)
	response := helper.BuildSuccessResponse("", poinSiswa)
	ctx.JSON(http.StatusOK, response)
	return
}

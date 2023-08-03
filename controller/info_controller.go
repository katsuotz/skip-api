package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
	"strconv"
)

type InfoController interface {
	CountPoinLog(ctx *gin.Context)
	MaxPoinSiswa(ctx *gin.Context)
	MinPoinSiswa(ctx *gin.Context)
	AvgPoinSiswa(ctx *gin.Context)
	ListPoinSiswa(ctx *gin.Context)
	ListPoinLog(ctx *gin.Context)
	ListCountPoinLog(ctx *gin.Context)
	GraphCountPoinLog(ctx *gin.Context)
	CountPoinSiswaTotal(ctx *gin.Context)
}

type infoController struct {
	PoinLogRepository   repository.PoinLogRepository
	PoinSiswaRepository repository.PoinSiswaRepository
}

func NewInfoController(
	poinLogRepository repository.PoinLogRepository,
	poinSiswaRepository repository.PoinSiswaRepository,
) InfoController {
	return &infoController{
		poinLogRepository,
		poinSiswaRepository,
	}
}

func (c *infoController) CountPoinLog(ctx *gin.Context) {
	poinType := ctx.DefaultQuery("type", "")
	kelasID := ctx.DefaultQuery("kelas_id", "")
	jurusanID := ctx.DefaultQuery("jurusan_id", "")
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	maxPoin := ctx.DefaultQuery("max_poin", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinLogRepository.CountPoinLog(ctx, poinType, kelasID, jurusanID, tahunAjarID, pegawaiID, maxPoin)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) MaxPoinSiswa(ctx *gin.Context) {
	kelasID := ctx.DefaultQuery("kelas_id", "")
	jurusanID := ctx.DefaultQuery("jurusan_id", "")
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinSiswaRepository.CountPoinSiswa(ctx, "max", kelasID, jurusanID, tahunAjarID, pegawaiID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) MinPoinSiswa(ctx *gin.Context) {
	kelasID := ctx.DefaultQuery("kelas_id", "")
	jurusanID := ctx.DefaultQuery("jurusan_id", "")
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinSiswaRepository.CountPoinSiswa(ctx, "min", kelasID, jurusanID, tahunAjarID, pegawaiID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) AvgPoinSiswa(ctx *gin.Context) {
	kelasID := ctx.DefaultQuery("kelas_id", "")
	jurusanID := ctx.DefaultQuery("jurusan_id", "")
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinSiswaRepository.CountPoinSiswa(ctx, "avg", kelasID, jurusanID, tahunAjarID, pegawaiID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) ListPoinSiswa(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	search := ctx.DefaultQuery("search", "")
	order := ctx.DefaultQuery("order", "asc")
	orderBy := ctx.DefaultQuery("order_by", "nama")
	maxPoin := ctx.DefaultQuery("max_poin", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinSiswaRepository.GetPoinSiswaPagination(ctx, pageInt, perPageInt, order, orderBy, search, tahunAjarID, pegawaiID, maxPoin)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) ListPoinLog(ctx *gin.Context) {
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

func (c *infoController) ListCountPoinLog(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	order := ctx.DefaultQuery("order", "asc")
	orderBy := ctx.DefaultQuery("order_by", "nama")
	groupBy := ctx.DefaultQuery("group_by", "siswa")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	poinType := ctx.DefaultQuery("type", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinLogRepository.GetCountPoinLogPagination(ctx, pageInt, perPageInt, order, orderBy, groupBy, tahunAjarID, poinType, pegawaiID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) GraphCountPoinLog(ctx *gin.Context) {
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	poinType := ctx.DefaultQuery("type", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinLogRepository.GetCountPoinLogPaginationByMonth(ctx, tahunAjarID, poinType, pegawaiID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) CountPoinSiswaTotal(ctx *gin.Context) {
	kelasID := ctx.DefaultQuery("kelas_id", "")
	jurusanID := ctx.DefaultQuery("jurusan_id", "")
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")
	maxPoin := ctx.DefaultQuery("max_poin", "")
	role := ctx.MustGet("role")
	pegawaiID := 0

	if role == "guru" {
		pegawaiID = int(ctx.MustGet("pegawai_id").(float64))
	}

	result := c.PoinSiswaRepository.CountPoinSiswaTotal(ctx, kelasID, jurusanID, tahunAjarID, pegawaiID, maxPoin)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

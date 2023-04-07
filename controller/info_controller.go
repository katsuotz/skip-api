package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
	"strconv"
)

type InfoController interface {
	CountPoin(ctx *gin.Context)
	ListPoin(ctx *gin.Context)
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

func (c *infoController) CountPoin(ctx *gin.Context) {
	poinType := ctx.DefaultQuery("type", "")
	kelasID := ctx.DefaultQuery("kelas_id", "")
	jurusanID := ctx.DefaultQuery("jurusan_id", "")

	result := c.PoinLogRepository.CountPoin(ctx, poinType, kelasID, jurusanID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *infoController) ListPoin(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	search := ctx.DefaultQuery("search", "")
	order := ctx.DefaultQuery("order", "asc")
	orderBy := ctx.DefaultQuery("order_by", "nama")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	tahunAjarID := ctx.DefaultQuery("tahun_ajar_id", "")

	result := c.PoinSiswaRepository.GetPoinSiswaPagination(ctx, pageInt, perPageInt, order, orderBy, search, tahunAjarID)
	response := helper.BuildSuccessResponse("", result)
	ctx.JSON(http.StatusOK, response)
	return
}

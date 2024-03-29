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

type DataPoinController interface {
	GetDataPoin(ctx *gin.Context)
	GetDataPoinByID(ctx *gin.Context)
	CreateDataPoin(ctx *gin.Context)
	UpdateDataPoin(ctx *gin.Context)
	DeleteDataPoin(ctx *gin.Context)
}

type dataPoinController struct {
	DataPoinRepository repository.DataPoinRepository
}

func NewDataPoinController(dataPoinRepository repository.DataPoinRepository) DataPoinController {
	return &dataPoinController{
		dataPoinRepository,
	}
}

func (c *dataPoinController) GetDataPoin(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	search := ctx.DefaultQuery("search", "")
	poinType := ctx.DefaultQuery("type", "")
	category := ctx.DefaultQuery("category", "")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	dataPoin := c.DataPoinRepository.GetDataPoin(ctx, pageInt, perPageInt, search, poinType, category)
	response := helper.BuildSuccessResponse("", dataPoin)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *dataPoinController) GetDataPoinByID(ctx *gin.Context) {
	dataPoinID, err := strconv.Atoi(ctx.Param("data_poin_id"))

	if err != nil || dataPoinID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	kelas := c.DataPoinRepository.GetDataPoinByID(ctx, dataPoinID)

	if kelas.ID == 0 {
		response := helper.BuildErrorResponse("Not Found", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response := helper.BuildSuccessResponse("", kelas)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *dataPoinController) CreateDataPoin(ctx *gin.Context) {
	req := dto.DataPoinRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	dataPoin := entity.DataPoin{
		Title:        req.Title,
		Description:  req.Description,
		Poin:         req.Poin,
		Type:         req.Type,
		Category:     req.Category,
		Penanganan:   req.Penanganan,
		TindakLanjut: req.TindakLanjut,
	}

	_, err := c.DataPoinRepository.CreateDataPoin(ctx, dataPoin)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Poin berhasil dibuat", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *dataPoinController) UpdateDataPoin(ctx *gin.Context) {
	req := dto.DataPoinRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	dataPoinID, err := strconv.Atoi(ctx.Param("data_poin_id"))
	if err != nil || dataPoinID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newDataPoin := entity.DataPoin{
		ID:          dataPoinID,
		Title:       req.Title,
		Description: req.Description,
		Poin:        req.Poin,
		Type:        req.Type,
		Category:    req.Category,
	}

	_, err = c.DataPoinRepository.UpdateDataPoin(ctx, newDataPoin)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Poin berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *dataPoinController) DeleteDataPoin(ctx *gin.Context) {
	dataPoinID, err := strconv.Atoi(ctx.Param("data_poin_id"))
	if err != nil || dataPoinID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.DataPoinRepository.DeleteDataPoin(ctx, dataPoinID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Poin berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

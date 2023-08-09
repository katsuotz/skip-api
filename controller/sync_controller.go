package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
	"strconv"
)

type SyncController interface {
	GetSync(ctx *gin.Context)
	Sync(ctx *gin.Context)
	SyncPassword(ctx *gin.Context)
}

type sitiController struct {
	SyncRepository repository.SyncRepository
}

func NewSyncController(sitiRepository repository.SyncRepository) SyncController {
	return &sitiController{
		sitiRepository,
	}
}

func (c *sitiController) GetSync(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	tahunAjar := c.SyncRepository.GetSync(ctx, pageInt, perPageInt)
	response := helper.BuildSuccessResponse("", tahunAjar)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *sitiController) Sync(ctx *gin.Context) {
	isOnProgress := c.SyncRepository.IsOnProgress(ctx, "siti")

	if isOnProgress == true {
		response := helper.BuildSuccessResponse("There's still synchronize in progress", nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	go c.SyncRepository.Sync(ctx)
	response := helper.BuildSuccessResponse("Synchronize Starting", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *sitiController) SyncPassword(ctx *gin.Context) {
	isOnProgress := c.SyncRepository.IsOnProgress(ctx, "password")

	if isOnProgress == true {
		response := helper.BuildSuccessResponse("There's still synchronize in progress", nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	c.SyncRepository.SyncPassword(ctx)
	response := helper.BuildSuccessResponse("Synchronize Password Done", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

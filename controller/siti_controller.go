package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
)

type SitiController interface {
	Sync(ctx *gin.Context)
	SyncPassword(ctx *gin.Context)
}

type sitiController struct {
	SitiRepository repository.SitiRepository
}

func NewSitiController(sitiRepository repository.SitiRepository) SitiController {
	return &sitiController{
		sitiRepository,
	}
}

func (c *sitiController) Sync(ctx *gin.Context) {
	c.SitiRepository.Sync(ctx)
	response := helper.BuildSuccessResponse("Synchronize Done", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *sitiController) SyncPassword(ctx *gin.Context) {
	c.SitiRepository.Sync(ctx)
	response := helper.BuildSuccessResponse("Synchronize Password Done", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

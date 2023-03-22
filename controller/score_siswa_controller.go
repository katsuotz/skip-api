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

type ScoreSiswaController interface {
	GetScoreSiswa(ctx *gin.Context)
	CreateScoreSiswa(ctx *gin.Context)
	UpdateScoreSiswa(ctx *gin.Context)
	DeleteScoreSiswa(ctx *gin.Context)
}

type scoreSiswaController struct {
	ScoreSiswaRepository repository.ScoreSiswaRepository
}

func NewScoreSiswaController(scoreSiswaRepository repository.ScoreSiswaRepository) ScoreSiswaController {
	return &scoreSiswaController{
		scoreSiswaRepository,
	}
}

func (c *scoreSiswaController) GetScoreSiswa(ctx *gin.Context) {
	scoreSiswa := c.ScoreSiswaRepository.GetScoreSiswa(ctx)
	response := helper.BuildSuccessResponse("", scoreSiswa)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *scoreSiswaController) CreateScoreSiswa(ctx *gin.Context) {
	req := dto.ScoreSiswaRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.ScoreSiswaRepository.AddScoreSiswa(ctx, req)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Score created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *scoreSiswaController) UpdateScoreSiswa(ctx *gin.Context) {
	req := dto.ScoreSiswaRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	scoreSiswaID, err := strconv.Atoi(ctx.Param("data_score_id"))
	if err != nil || scoreSiswaID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newScoreSiswa := entity.ScoreSiswa{
		ID: scoreSiswaID,
		//Title:       req.Title,
		//Description: req.Description,
		Score: req.Score,
	}

	_, err = c.ScoreSiswaRepository.UpdateScoreSiswa(ctx, newScoreSiswa)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Score updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *scoreSiswaController) DeleteScoreSiswa(ctx *gin.Context) {
	scoreSiswaID, err := strconv.Atoi(ctx.Param("data_score_id"))
	if err != nil || scoreSiswaID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.ScoreSiswaRepository.DeleteScoreSiswa(ctx, scoreSiswaID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Score deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

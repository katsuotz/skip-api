package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/service"
	"net/http"
	"strconv"
)

type DataScoreController interface {
	GetDataScore(ctx *gin.Context)
	CreateDataScore(ctx *gin.Context)
	UpdateDataScore(ctx *gin.Context)
	DeleteDataScore(ctx *gin.Context)
}

type dataScoreController struct {
	DataScoreRepository repository.DataScoreRepository
	JWTService          service.JWTService
}

func NewDataScoreController(dataScoreRepository repository.DataScoreRepository, jwtService service.JWTService) DataScoreController {
	return &dataScoreController{
		dataScoreRepository,
		jwtService,
	}
}

func (c *dataScoreController) GetDataScore(ctx *gin.Context) {
	dataScore := c.DataScoreRepository.GetDataScore(ctx)
	response := helper.BuildSuccessResponse("", dataScore)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *dataScoreController) CreateDataScore(ctx *gin.Context) {
	req := dto.DataScoreRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	dataScore := entity.DataScore{
		Title:       req.Title,
		Description: req.Description,
		Score:       req.Score,
	}

	_, err := c.DataScoreRepository.CreateDataScore(ctx, dataScore)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Score created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *dataScoreController) UpdateDataScore(ctx *gin.Context) {
	req := dto.DataScoreRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	dataScoreID, err := strconv.Atoi(ctx.Param("data_score_id"))
	if err != nil || dataScoreID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newDataScore := entity.DataScore{
		ID:          dataScoreID,
		Title:       req.Title,
		Description: req.Description,
		Score:       req.Score,
	}

	_, err = c.DataScoreRepository.UpdateDataScore(ctx, newDataScore)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Score updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *dataScoreController) DeleteDataScore(ctx *gin.Context) {
	dataScoreID, err := strconv.Atoi(ctx.Param("data_score_id"))
	if err != nil || dataScoreID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.DataScoreRepository.DeleteDataScore(ctx, dataScoreID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Data Score deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

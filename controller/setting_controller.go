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

type SettingController interface {
	GetSetting(ctx *gin.Context)
	CreateSetting(ctx *gin.Context)
	UpdateSetting(ctx *gin.Context)
	DeleteSetting(ctx *gin.Context)
}

type settingController struct {
	SettingRepository repository.SettingRepository
}

func NewSettingController(settingRepository repository.SettingRepository) SettingController {
	return &settingController{
		settingRepository,
	}
}

func (c *settingController) GetSetting(ctx *gin.Context) {
	key := ctx.DefaultQuery("key", "")

	if key != "" {
		setting := c.SettingRepository.GetSettingByKey(ctx, key)

		if setting.ID == 0 {
			response := helper.BuildErrorResponse("Data not found", nil, nil)
			ctx.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		}

		response := helper.BuildSuccessResponse("", setting)
		ctx.JSON(http.StatusOK, response)
		return
	}

	setting := c.SettingRepository.GetSetting(ctx)
	response := helper.BuildSuccessResponse("", setting)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *settingController) CreateSetting(ctx *gin.Context) {
	req := dto.SettingRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	setting := entity.Setting{
		Key:   req.Key,
		Value: req.Value,
	}

	_, err := c.SettingRepository.CreateSetting(ctx, setting)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Setting berhasil dibuat", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *settingController) UpdateSetting(ctx *gin.Context) {
	req := dto.SettingRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	settingID, err := strconv.Atoi(ctx.Param("setting_id"))
	if err != nil || settingID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newSetting := entity.Setting{
		ID:    settingID,
		Key:   req.Key,
		Value: req.Value,
	}

	_, err = c.SettingRepository.UpdateSetting(ctx, newSetting)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Setting berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *settingController) DeleteSetting(ctx *gin.Context) {
	settingID, err := strconv.Atoi(ctx.Param("setting_id"))
	if err != nil || settingID == 0 {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.SettingRepository.DeleteSetting(ctx, settingID)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Setting berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"net/http"
)

type ProfileController interface {
	GetMyProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
}

type profileController struct {
	ProfileRepository repository.ProfileRepository
}

func NewProfileController(profileRepository repository.ProfileRepository) ProfileController {
	return &profileController{
		profileRepository,
	}
}

func (c *profileController) GetMyProfile(ctx *gin.Context) {
	userID := int(ctx.MustGet("user_id").(float64))

	profile := c.ProfileRepository.FindProfileWithJoinByID(ctx, userID)
	response := helper.BuildSuccessResponse("", profile)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *profileController) UpdateProfile(ctx *gin.Context) {
	req := dto.ProfileRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userID := int(ctx.MustGet("user_id").(float64))

	profile := c.ProfileRepository.FindProfileByID(ctx, userID)
	if profile.ID == 0 {
		response := helper.BuildErrorResponse("Unauthorized", nil, nil)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	tanggalLahir, _ := helper.StringToDate(req.TanggalLahir)

	newProfile := entity.Profile{
		ID:           profile.ID,
		Nama:         req.Nama,
		JenisKelamin: req.JenisKelamin,
		TempatLahir:  req.TempatLahir,
		TanggalLahir: tanggalLahir,
		Foto:         req.Foto,
	}

	_, err := c.ProfileRepository.UpdateProfile(ctx, newProfile)

	if err != nil {
		response := helper.BuildErrorResponse("Gagal memproses permintaan", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Profile berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

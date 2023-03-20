package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gitlab.com/katsuotz/skip-api/helper"
	"gitlab.com/katsuotz/skip-api/repository"
	"gitlab.com/katsuotz/skip-api/service"
	"net/http"
)

type ProfileController interface {
	UpdateProfile(ctx *gin.Context)
}

type profileController struct {
	ProfileRepository repository.ProfileRepository
	JWTService        service.JWTService
}

func NewProfileController(profileRepository repository.ProfileRepository, jwtService service.JWTService) ProfileController {
	return &profileController{
		profileRepository,
		jwtService,
	}
}

func (c *profileController) UpdateProfile(ctx *gin.Context) {
	req := dto.ProfileRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userID := int(ctx.MustGet("user_id").(float64))

	profile := c.ProfileRepository.FindByProfileID(ctx, userID)
	if profile.ProfileID == 0 {
		response := helper.BuildErrorResponse("Unauthorized", nil, nil)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	newProfile := entity.Profile{
		ProfileID:    profile.ProfileID,
		Nama:         req.Nama,
		JenisKelamin: req.JenisKelamin,
		TempatLahir:  req.TempatLahir,
		TanggalLahir: req.TanggalLahir,
	}

	_, err := c.ProfileRepository.UpdateProfile(ctx, newProfile)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Profile updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

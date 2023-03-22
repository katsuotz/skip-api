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

type KelasController interface {
	GetKelas(ctx *gin.Context)
	CreateKelas(ctx *gin.Context)
	UpdateKelas(ctx *gin.Context)
	DeleteKelas(ctx *gin.Context)
	AddSiswaToKelas(ctx *gin.Context)
	RemoveSiswaFromKelas(ctx *gin.Context)
}

type kelasController struct {
	KelasRepository repository.KelasRepository
	JWTService      service.JWTService
}

func NewKelasController(kelasRepository repository.KelasRepository, jwtService service.JWTService) KelasController {
	return &kelasController{
		kelasRepository,
		jwtService,
	}
}

func (c *kelasController) GetKelas(ctx *gin.Context) {
	jurusanID := ctx.DefaultQuery("jurusan_id", "")
	tahunAjarID := ctx.Query("tahun_ajar_id")

	if tahunAjarID == "" {
		response := helper.BuildErrorResponse("Tahun Ajar Needed", nil, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	kelas := c.KelasRepository.GetKelas(ctx, jurusanID, tahunAjarID)
	response := helper.BuildSuccessResponse("", kelas)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *kelasController) CreateKelas(ctx *gin.Context) {
	req := dto.KelasRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	kelas := entity.Kelas{
		NamaKelas:   req.NamaKelas,
		JurusanID:   req.JurusanID,
		TahunAjarID: req.TahunAjarID,
		GuruID:      req.GuruID,
	}

	_, err := c.KelasRepository.CreateKelas(ctx, kelas)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Kelas created successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *kelasController) UpdateKelas(ctx *gin.Context) {
	req := dto.KelasRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	kelasID, err := strconv.Atoi(ctx.Param("kelas_id"))
	if err != nil || kelasID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newKelas := entity.Kelas{
		ID:          kelasID,
		NamaKelas:   req.NamaKelas,
		JurusanID:   req.JurusanID,
		TahunAjarID: req.TahunAjarID,
	}

	_, err = c.KelasRepository.UpdateKelas(ctx, newKelas)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Kelas updated successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *kelasController) DeleteKelas(ctx *gin.Context) {
	kelasID, err := strconv.Atoi(ctx.Param("kelas_id"))
	if err != nil || kelasID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.KelasRepository.DeleteKelas(ctx, int(kelasID))

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Kelas deleted successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *kelasController) AddSiswaToKelas(ctx *gin.Context) {
	req := dto.DetailKelasRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	kelasID, err := strconv.Atoi(ctx.Param("kelas_id"))
	if err != nil || kelasID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.KelasRepository.AddSiswaToKelas(ctx, kelasID, req.SiswaID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa added successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *kelasController) RemoveSiswaFromKelas(ctx *gin.Context) {
	req := dto.DetailKelasRequest{}
	errDTO := ctx.ShouldBindJSON(&req)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	kelasID, err := strconv.Atoi(ctx.Param("kelas_id"))
	if err != nil || kelasID == 0 {
		response := helper.BuildErrorResponse("Failed to process request", nil, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = c.KelasRepository.RemoveSiswaFromKelas(ctx, kelasID, req.SiswaID)

	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse("Siswa removed successfully", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

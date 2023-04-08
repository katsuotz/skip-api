package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/helper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileController interface {
	Upload(ctx *gin.Context)
}

type fileController struct {
}

func NewFileController() FileController {
	return &fileController{}
}

func (c *fileController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var formData dto.File

	err = ctx.ShouldBind(&formData)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		response := helper.BuildErrorResponse("Internal server error", err, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	defer uploadedFile.Close()

	fileContent, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		response := helper.BuildErrorResponse("Internal server error", err, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)

	split := strings.Split(file.Filename, ".")
	ext := split[len(split)-1]

	filename := formData.Filename + "-" + timestamp + helper.RandomString(4) + "." + ext

	dir := "storage"

	if formData.Folder != "" {
		dir += "/" + formData.Folder
	}

	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	destinationFilename := filepath.Join(dir, filename)
	destinationFile, err := os.Create(destinationFilename)
	if err != nil {
		response := helper.BuildErrorResponse("Internal server error", err, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	defer destinationFile.Close()

	_, err = destinationFile.Write(fileContent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := helper.BuildSuccessResponse("File uploaded successfully", gin.H{
		"url": dir + "/" + filename,
	})
	ctx.JSON(http.StatusOK, response)
}

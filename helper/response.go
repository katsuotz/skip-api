package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}
	return fe.Tag()
}

func BuildSuccessResponse(message string, data interface{}) any {
	res := gin.H{
		"status":  true,
		"message": message,
	}

	if data != nil {
		res["data"] = data
	}

	return res
}

func BuildErrorResponse(message string, err error, data interface{}) any {
	res := gin.H{
		"status":  false,
		"message": message,
	}
	var ve validator.ValidationErrors
	if err != nil && errors.As(err, &ve) {
		errMsg := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			errMsg[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
		}
		res["error"] = errMsg
	}

	if data != nil {
		res["data"] = data
	}

	return res
}

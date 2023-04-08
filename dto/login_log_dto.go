package dto

import (
	"gitlab.com/katsuotz/skip-api/entity"
)

type LoginLogResponse struct {
	Username string `json:"username"`
	Nama     string `json:"nama"`
	Role     string `json:"role"`
	entity.LoginLog
}

type LoginLogPagination struct {
	Data       []LoginLogResponse `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

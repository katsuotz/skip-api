package dto

import "gitlab.com/katsuotz/skip-api/entity"

type SiswaRequest struct {
	Nis string `json:"nis" binding:"required" form:"nis"`
	ProfileRequest
}

type SiswaResponse struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	entity.Siswa
	entity.Profile
}

type SiswaPagination struct {
	Data       []SiswaResponse `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

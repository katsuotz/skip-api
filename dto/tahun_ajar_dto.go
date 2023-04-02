package dto

import "gitlab.com/katsuotz/skip-api/entity"

type TahunAjarRequest struct {
	TahunAjar string `json:"tahun_ajar" binding:"required" form:"tahun_ajar"`
}

type TahunAjarPagination struct {
	Data       []entity.TahunAjar `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

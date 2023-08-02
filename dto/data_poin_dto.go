package dto

import "gitlab.com/katsuotz/skip-api/entity"

type DataPoinRequest struct {
	Title        string  `json:"title" binding:"required" form:"title"`
	Description  string  `json:"description" binding:"required" form:"description"`
	Poin         float64 `json:"poin" binding:"required" form:"poin"`
	Type         string  `json:"type" binding:"required" form:"type"`
	Category     string  `json:"category" binding:"required" form:"category"`
	Penanganan   string  `json:"penanganan" form:"penanganan"`
	TindakLanjut string  `json:"tindak_lanjut" form:"tindak_lanjut"`
}

type DataPoinPagination struct {
	Data       []entity.DataPoin `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

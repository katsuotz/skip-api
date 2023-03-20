package dto

type TahunAjarRequest struct {
	TahunAjar string `json:"tahun_ajar" binding:"required" form:"tahun_ajar"`
}

package dto

type PoinSiswaRequest struct {
	Title        string  `json:"title" binding:"required" form:"title"`
	Description  string  `json:"description" binding:"required" form:"description"`
	Poin         float64 `json:"poin" binding:"required" form:"poin"`
	Type         string  `json:"type" binding:"required" form:"type"`
	GuruID       int     `json:"guru_id" form:"guru_id"`
	SiswaKelasID int     `json:"siswa_kelas_id" binding:"required" form:"siswa_kelas_id"`
}

type UpdatePoinLogRequest struct {
	Description string `json:"description" binding:"required" form:"description"`
}

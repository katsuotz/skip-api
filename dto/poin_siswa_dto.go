package dto

import "gitlab.com/katsuotz/skip-api/entity"

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

type PoinSiswaResponse struct {
	Nis       string  `json:"nis"`
	Nama      string  `json:"nama"`
	NamaKelas string  `json:"nama_kelas"`
	Poin      float64 `json:"poin"`
	entity.Base
}

type PoinLogResponse struct {
	ID          int     `json:"id"`
	Poin        float64 `json:"poin"`
	PoinBefore  float64 `json:"poin_before"`
	PoinAfter   float64 `json:"poin_after"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	GuruID      int     `json:"guru_id"`
	Nip         string  `json:"nip"`
	NamaGuru    string  `json:"nama_guru"`
	entity.Base
}

type PoinLogPagination struct {
	Data       []PoinLogResponse `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

type PoinLogSiswaByKelas struct {
	Kelas entity.Kelas      `json:"kelas"`
	Data  []PoinLogResponse `json:"data"`
}

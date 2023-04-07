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

type PoinKelasResponse struct {
	NamaKelas string  `json:"nama_kelas"`
	TahunAjar string  `json:"tahun_ajar"`
	Poin      float64 `json:"poin"`
}

type PoinJurusanResponse struct {
	NamaJurusan string  `json:"nama_jurusan"`
	TahunAjar   string  `json:"tahun_ajar"`
	Poin        float64 `json:"poin"`
}

type PoinSiswaPagination struct {
	Data       []PoinSiswaResponse `json:"data"`
	Pagination Pagination          `json:"pagination"`
}

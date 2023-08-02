package dto

import "gitlab.com/katsuotz/skip-api/entity"

type PoinSiswaRequest struct {
	Title        string  `json:"title" binding:"required" form:"title"`
	Description  string  `json:"description" binding:"required" form:"description"`
	Penanganan   string  `json:"penanganan" form:"penanganan"`
	Poin         float64 `json:"poin" binding:"required" form:"poin"`
	Type         string  `json:"type" binding:"required" form:"type"`
	PegawaiID    int     `json:"pegawai_id" form:"pegawai_id"`
	File         string  `json:"file" form:"file"`
	SiswaKelasID int     `json:"siswa_kelas_id" binding:"required" form:"siswa_kelas_id"`
	DataPoinID   int     `json:"data_poin_id" binding:"required" form:"data_poin_id"`
}

type UpdatePoinLogRequest struct {
	TindakLanjut string `json:"tindak_lanjut" binding:"required" form:"tindak_lanjut"`
}

type PoinSiswaResponse struct {
	Nis       string  `json:"nis"`
	Nama      string  `json:"nama"`
	NamaKelas string  `json:"nama_kelas"`
	Foto      string  `json:"foto"`
	Poin      float64 `json:"poin"`
	entity.Base
}

type PoinKelasResponse struct {
	NamaKelas string  `json:"nama_kelas"`
	Poin      float64 `json:"poin"`
}

type PoinJurusanResponse struct {
	NamaJurusan string  `json:"nama_jurusan"`
	TahunAjar   string  `json:"tahun_ajar"`
	Poin        float64 `json:"poin"`
}

type PoinJurusanWithKelasResponse struct {
	Jurusan PoinJurusanResponse `json:"jurusan"`
	Data    []PoinKelasResponse `json:"data"`
}

type PoinSiswaPagination struct {
	Data       []PoinSiswaResponse `json:"data"`
	Pagination Pagination          `json:"pagination"`
}

package dto

import "gitlab.com/katsuotz/skip-api/entity"

type KelasRequest struct {
	NamaKelas   string `json:"nama_kelas" binding:"required" form:"nama_kelas"`
	JurusanID   int    `json:"jurusan_id" binding:"required" form:"jurusan_id"`
	TahunAjarID int    `json:"tahun_ajar_id" binding:"required" form:"tahun_ajar_id"`
	GuruID      int    `json:"guru_id" binding:"required" form:"guru_id"`
}

type SiswaKelasRequest struct {
	SiswaID []int `json:"siswa_id" binding:"required" form:"siswa_id"`
}

type KelasResponse struct {
	entity.Kelas
	entity.Profile
	entity.Guru
}

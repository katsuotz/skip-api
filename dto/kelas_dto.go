package dto

import "gitlab.com/katsuotz/skip-api/entity"

type KelasRequest struct {
	NamaKelas   string `json:"nama_kelas" binding:"required" form:"nama_kelas"`
	JurusanID   int    `json:"jurusan_id" binding:"required" form:"jurusan_id"`
	TahunAjarID int    `json:"tahun_ajar_id" binding:"required" form:"tahun_ajar_id"`
	PegawaiID   int    `json:"pegawai_id" binding:"required" form:"pegawai_id"`
	Tingkat     int    `json:"tingkat" binding:"required" form:"tingkat"`
}

type SiswaKelasRequest struct {
	SiswaID []int `json:"siswa_id" binding:"required" form:"siswa_id"`
}

type SiswaNaikKelasRequest struct {
	SiswaKelasID []int `json:"siswa_kelas_id" binding:"required" form:"siswa_kelas_id"`
}

type KelasResponse struct {
	ID int `json:"id"`
	entity.Kelas
	entity.Profile
	entity.Pegawai
}

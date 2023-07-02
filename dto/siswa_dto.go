package dto

import "gitlab.com/katsuotz/skip-api/entity"

type SiswaRequest struct {
	Nis string `json:"nis" binding:"required" form:"nis"`
	ProfileRequest
}

type SiswaResponse struct {
	ID               int      `json:"id"`
	UserID           int      `json:"user_id"`
	Poin             *float64 `json:"poin,omitempty"`
	SiswaKelasID     *int     `json:"siswa_kelas_id,omitempty"`
	TotalPenghargaan *float64 `json:"total_penghargaan,omitempty"`
	TotalPelanggaran *float64 `json:"total_pelanggaran,omitempty"`
	NamaKelas        *string  `json:"nama_kelas,omitempty"`
	entity.Siswa
	entity.Profile
}

type SiswaPagination struct {
	Data       []SiswaResponse `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type SiswaDetailLog struct {
	Siswa SiswaResponse         `json:"siswa"`
	Log   []PoinLogSiswaByKelas `json:"log"`
}

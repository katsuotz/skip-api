package dto

import "gitlab.com/katsuotz/skip-api/entity"

type ProfileRequest struct {
	Nama         string `json:"nama" binding:"required" form:"nama"`
	JenisKelamin string `json:"jenis_kelamin" binding:"required" form:"jenis_kelamin"`
	TempatLahir  string `json:"tempat_lahir" binding:"required" form:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required" form:"tanggal_lahir"`
}

type ProfileResponse struct {
	ID       int     `json:"id"`
	Nip      *string `json:"nip,omitempty"`
	Nis      *string `json:"nis,omitempty"`
	TipeGuru *string `json:"tipe_guru,omitempty"`
	entity.Profile
	entity.Guru
	entity.Siswa
	entity.User
}

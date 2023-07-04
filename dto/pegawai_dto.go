package dto

import "gitlab.com/katsuotz/skip-api/entity"

type PegawaiRequest struct {
	Nip         string `json:"nip" binding:"required" form:"nip"`
	TipePegawai string `json:"tipe_pegawai" binding:"required" form:"tipe_pegawai"`
	Username    string `json:"username" binding:"required" form:"username"`
	Password    string `json:"password" form:"password"`
	ProfileRequest
}

type CreatePegawaiRequest struct {
	Password string `json:"password" binding:"required" form:"password"`
	PegawaiRequest
}

type PegawaiResponse struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	entity.Pegawai
	entity.Profile
}

type PegawaiPagination struct {
	Data       []PegawaiResponse `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

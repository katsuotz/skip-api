package dto

import "gitlab.com/katsuotz/skip-api/entity"

type GuruRequest struct {
	Nip      string `json:"nip" binding:"required" form:"nip"`
	TipeGuru string `json:"tipe_guru" binding:"required" form:"tipe_guru"`
	Username string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
	ProfileRequest
}

type GuruResponse struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	entity.Guru
	entity.Profile
}

type GuruPagination struct {
	Data       []GuruResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

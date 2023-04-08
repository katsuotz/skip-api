package dto

import "gitlab.com/katsuotz/skip-api/entity"

type LoginRequest struct {
	Username string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID       int     `json:"-"`
	Nip      *string `json:"nip,omitempty"`
	Nis      *string `json:"nis,omitempty"`
	GuruID   *string `json:"-"`
	SiswaID  *string `json:"-"`
	TipeGuru *string `json:"tipe_guru,omitempty"`
	entity.Profile
	entity.Guru
	entity.Siswa
	entity.User
}

type UpdatePasswordRequest struct {
	OldPassword          string `json:"old_password" binding:"required" form:"old_password"`
	Password             string `json:"password" binding:"required" form:"password"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required" form:"password_confirmation"`
}

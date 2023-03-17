package dto

import "gitlab.com/katsuotz/skip-api/entity"

type LoginRequest struct {
	Username string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

type LoginResponse struct {
	User  entity.User `json:"user"`
	Token string      `json:"token"`
}

package dto

type CreateSiswaRequest struct {
	Nis string `json:"nis" binding:"required" form:"nis"`
	ProfileRequest
}

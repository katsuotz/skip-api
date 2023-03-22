package dto

type SettingRequest struct {
	Key   string `json:"key" binding:"required" form:"key"`
	Value string `json:"value" binding:"required" form:"value"`
}

package dto

type File struct {
	Filename string `json:"filename" binding:"required" form:"filename"`
	Folder   string `json:"folder" binding:"required" form:"folder"`
}

package dto

type DataScoreRequest struct {
	Title       string  `json:"title" binding:"required" form:"title"`
	Description string  `json:"description" binding:"required" form:"description"`
	Score       float64 `json:"score" binding:"required" form:"score"`
}

package dto

type ScoreSiswaRequest struct {
	Title        string  `json:"title" binding:"required" form:"title"`
	Description  string  `json:"description" binding:"required" form:"description"`
	Score        float64 `json:"score" binding:"required" form:"score"`
	Type         string  `json:"type" binding:"required" form:"type"`
	GuruID       int     `json:"guru_id" binding:"required" form:"guru_id"`
	SiswaKelasID int     `json:"siswa_kelas_id" binding:"required" form:"siswa_kelas_id"`
}

type UpdateScoreLogRequest struct {
	Description string `json:"description" binding:"required" form:"description"`
}

package dto

type JurusanRequest struct {
	NamaJurusan string `json:"nama_jurusan" binding:"required" form:"nama_jurusan"`
}

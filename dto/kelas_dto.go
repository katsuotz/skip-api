package dto

type KelasRequest struct {
	NamaKelas   string `json:"nama_kelas" binding:"required" form:"nama_kelas"`
	JurusanID   int    `json:"jurusan_id" binding:"required" form:"jurusan_id"`
	TahunAjarID int    `json:"tahun_ajar_id" binding:"required" form:"tahun_ajar_id"`
	GuruID      int    `json:"guru_id" binding:"required" form:"guru_id"`
}

package dto

type ProfileRequest struct {
	Nama         string `json:"nama" binding:"required" form:"nama"`
	JenisKelamin string `json:"jenis_kelamin" binding:"required" form:"jenis_kelamin"`
	TempatLahir  string `json:"tempat_lahir" binding:"required" form:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required" form:"tanggal_lahir"`
}

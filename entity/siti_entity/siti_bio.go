package siti_entity

import "time"

type Bio struct {
	IDBio        int       `json:"id_bio"`
	NamaLengkap  string    `json:"nama_lengkap"`
	JenisKelamin string    `json:"jenis_kelamin"`
	TempatLahir  string    `json:"tempat_lahir"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
}

func (Bio) TableName() string {
	return "t_bio"
}

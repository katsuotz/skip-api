package entity

import "time"

type Profile struct {
	ID           int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Nama         string    `gorm:"type:varchar(100)" json:"nama"`
	JenisKelamin string    `gorm:"type:varchar(10)" json:"jenis_kelamin"`
	TempatLahir  string    `gorm:"type:varchar(20)" json:"tempat_lahir"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	UserID       int       `json:"user_id"`
	Base
}

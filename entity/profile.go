package entity

import "time"

type Profile struct {
	ID           int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Nama         string    `gorm:"type:varchar(100)" json:"nama"`
	JenisKelamin string    `gorm:"type:varchar(10)" json:"jenis_kelamin,omitempty"`
	TempatLahir  string    `gorm:"type:varchar(20)" json:"tempat_lahir,omitempty"`
	TanggalLahir time.Time `gorm:"type:date" json:"tanggal_lahir,omitempty"`
	UserID       int       `gorm:"uniqueIndex" json:"user_id"`
	Base
}

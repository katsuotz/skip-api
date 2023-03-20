package entity

type Siswa struct {
	SiswaID int    `gorm:"primary_key;AUTO_INCREMENT" json:"siswa_id"`
	NIS     string `gorm:"type:varchar(20)" json:"nis"`
	UserID  int    `json:"user_id"`
	Base
}

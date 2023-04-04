package entity

type PoinSiswa struct {
	ID           int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	SiswaKelasID int     `gorm:"index" json:"siswa_kelas_id"`
	Poin         float64 `json:"poin"`
	Base
}

func (PoinSiswa) TableName() string {
	return "poin_siswa"
}

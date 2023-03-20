package entity

type DetailKelas struct {
	ID      int   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	KelasID int   `json:"kelas_id"`
	SiswaID int   `json:"siswa_id"`
	GuruID  int   `json:"guru_id,omitempty"`
	Kelas   Kelas `gorm:"foreignKey:id" json:"kelas"`
	Siswa   Siswa `gorm:"foreignKey:id" json:"siswa"`
	Base
}

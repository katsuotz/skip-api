package entity

type DetailKelas struct {
	ID      int   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	KelasID int   `gorm:"index" json:"kelas_id"`
	SiswaID int   `gorm:"index" json:"siswa_id"`
	Kelas   Kelas `gorm:"foreignKey:KelasID" json:"kelas"`
	Siswa   Siswa `gorm:"foreignKey:SiswaID" json:"siswa"`
	Base
}

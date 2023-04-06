package entity

type SiswaKelas struct {
	ID      int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	KelasID int    `gorm:"index" json:"kelas_id"`
	SiswaID int    `gorm:"index" json:"siswa_id"`
	Kelas   *Kelas `gorm:"foreignKey:KelasID" json:"kelas,omitempty"`
	Siswa   *Siswa `gorm:"foreignKey:SiswaID" json:"siswa,omitempty"`
	Base
}

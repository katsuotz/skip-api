package entity

type DetailKelas struct {
	DetailKelasID int   `gorm:"primary_key;AUTO_INCREMENT" json:"detail_kelas_id"`
	KelasID       int   `json:"kelas_id"`
	SiswaID       int   `json:"siswa_id"`
	GuruID        int   `json:"guru_id,omitempty"`
	Kelas         Kelas `gorm:"foreignKey:kelas_id" json:"kelas"`
	Siswa         Siswa `gorm:"foreignKey:siswa_id" json:"siswa"`
	Base
}

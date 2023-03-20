package entity

type Kelas struct {
	KelasID     int       `gorm:"primary_key;AUTO_INCREMENT" json:"kelas_id"`
	NamaKelas   string    `gorm:"type:varchar(20)" json:"nama_kelas"`
	JurusanID   int       `json:"jurusan_id"`
	TahunAjarID int       `json:"tahun_ajar_id"`
	GuruID      int       `json:"guru_id"`
	Jurusan     Jurusan   `gorm:"foreignKey:jurusan_id" json:"jurusan"`
	TahunAjar   TahunAjar `gorm:"foreignKey:tahun_ajar_id" json:"tahun_ajar"`
	Guru        Guru      `gorm:"foreignKey:guru_id" json:"guru"`
	Base
}

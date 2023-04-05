package entity

type Kelas struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	NamaKelas   string `gorm:"type:varchar(20)" json:"nama_kelas"`
	JurusanID   int    `json:"jurusan_id"`
	TahunAjarID int    `json:"tahun_ajar_id"`
	GuruID      int    `gorm:"index" json:"guru_id"`
	Guru        *Guru  `gorm:"foreignKey:guru_id" json:"guru,omitempty"`
	Base
}

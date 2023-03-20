package entity

type Jurusan struct {
	JurusanID   int    `gorm:"primary_key;AUTO_INCREMENT" json:"jurusan_id"`
	NamaJurusan string `gorm:"type:varchar(20)" json:"nama_jurusan"`
	Base
}

package entity

type Jurusan struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	NamaJurusan string `gorm:"type:varchar(50)" json:"nama_jurusan"`
	Base
}

func (Jurusan) TableName() string {
	return "jurusan"
}

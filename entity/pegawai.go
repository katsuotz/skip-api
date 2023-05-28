package entity

type Pegawai struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Nip         string `gorm:"type:varchar(20)" json:"nip"`
	TipePegawai string `gorm:"type:varchar(15)" json:"tipe_pegawai"`
	UserID      int    `gorm:"uniqueIndex" json:"user_id"`
	Base
}

func (Pegawai) TableName() string {
	return "pegawai"
}

package entity

type Guru struct {
	GuruID   int    `gorm:"primary_key;AUTO_INCREMENT" json:"guru_id"`
	NIP      string `gorm:"type:varchar(20)" json:"nip"`
	TipeGuru int    `gorm:"type:varchar(15)" json:"tipe_guru"`
	UserID   int    `json:"user_id"`
	Base
}

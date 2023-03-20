package entity

type Guru struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Nip      string `gorm:"type:varchar(20)" json:"nip"`
	TipeGuru string `gorm:"type:varchar(15)" json:"tipe_guru"`
	UserID   int    `json:"user_id"`
	Base
}

func (Guru) TableName() string {
	return "guru"
}

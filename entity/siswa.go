package entity

type Siswa struct {
	ID     int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Nis    string `gorm:"type:varchar(20)" json:"nis"`
	UserID int    `json:"user_id"`
	Base
}

func (Siswa) TableName() string {
	return "siswa"
}

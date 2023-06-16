package entity

type TahunAjar struct {
	ID        int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TahunAjar string `gorm:"type:varchar(10)" json:"tahun_ajar"`
	IsActive  bool   `gorm:"index;default:false" json:"is_active"`
	Base
}

func (TahunAjar) TableName() string {
	return "tahun_ajar"
}

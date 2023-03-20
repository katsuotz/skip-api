package entity

type TahunAjar struct {
	TahunAjarID int    `gorm:"primary_key;AUTO_INCREMENT" json:"tahun_ajar_id"`
	TahunAjar   string `gorm:"type:varchar(10)" json:"tahun_ajar"`
	IsActive    bool   `gorm:"default:false" json:"is_active"`
	Base
}

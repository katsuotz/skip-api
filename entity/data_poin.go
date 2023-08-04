package entity

type DataPoin struct {
	ID           int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title        string  `gorm:"type:varchar(200)" json:"title"`
	Description  string  `gorm:"type:text" json:"description"`
	Poin         float64 `json:"poin"`
	Penanganan   string  `gorm:"type:text" json:"penanganan,omitempty"`
	TindakLanjut string  `gorm:"type:text" json:"tindak_lanjut,omitempty"`
	Type         string  `gorm:"index;index:type_category_idx;type:varchar(25)" json:"type"`
	Category     string  `gorm:"index;index:type_category_idx;type:varchar(15)" json:"category"`
	Base
}

func (DataPoin) TableName() string {
	return "data_poin"
}

package entity

type DataPoin struct {
	ID          int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title       string  `gorm:"type:varchar(50)" json:"title"`
	Description string  `gorm:"type:text" json:"description"`
	Poin        float64 `json:"poin"`
	Type        string  `gorm:"type:varchar(25)" json:"type"`
	Category    string  `gorm:"type:varchar(15)" json:"category"`
	Base
}

func (DataPoin) TableName() string {
	return "data_poin"
}

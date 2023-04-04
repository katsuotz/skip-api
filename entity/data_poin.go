package entity

type DataPoin struct {
	ID          int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title       string  `gorm:"type:varchar(50)" json:"title"`
	Description string  `gorm:"type:text" json:"description"`
	Poin        float64 `json:"poin"`
	Type        string  `gorm:"type:varchar(20)" json:"type"`
	Base
}

func (DataPoin) TableName() string {
	return "data_poin"
}

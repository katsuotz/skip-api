package entity

type DataScore struct {
	ID          int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title       string  `gorm:"type:varchar(50)" json:"title"`
	Description string  `gorm:"type:text" json:"description"`
	Score       float64 `json:"score"`
	Base
}

func (DataScore) TableName() string {
	return "data_score"
}

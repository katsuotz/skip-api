package entity

type Sync struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Type        string `gorm:"type:varchar(30)" json:"type"`
	Status      string `gorm:"type:varchar(30)" json:"status"`
	Description string `json:"description"`
	Base
}

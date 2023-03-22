package entity

type Setting struct {
	ID        int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Key       string `gorm:"uniqueIndex" json:"key"`
	Value     string `gorm:"type:text" json:"value"`
	CanDelete bool   `gorm:"default:false" json:"can_delete"`
	Base
}

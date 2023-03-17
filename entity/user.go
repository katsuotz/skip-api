package entity

type User struct {
	UserID   int    `gorm:"primary_key;AUTO_INCREMENT" json:"user_id"`
	Username string `gorm:"type:varchar(100)" json:"username,omitempty"`
	Password string `gorm:"type:text" json:"-"`
	Role     string `gorm:"type:varchar(20)" json:"role"`
	Base
}

package entity

type User struct {
	ID       int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Username string  `gorm:"type:varchar(100)" json:"username,omitempty"`
	Password string  `gorm:"type:text" json:"-"`
	Role     string  `gorm:"type:varchar(20)" json:"role"`
	Profile  Profile `gorm:"references:id" json:"profile"`
	Base
}

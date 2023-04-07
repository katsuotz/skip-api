package entity

import "time"

type User struct {
	ID       int      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Username string   `gorm:"type:varchar(100)" json:"username,omitempty"`
	Password string   `gorm:"type:text" json:"-"`
	Role     string   `gorm:"type:varchar(20)" json:"role"`
	Profile  *Profile `gorm:"references:id" json:"profile,omitempty"`
	Base
}

type LoginLog struct {
	UserID    int        `gorm:"index" json:"user_id"`
	Action    string     `gorm:"type:varchar(30)" json:"action"`
	UserAgent string     `gorm:"type:varchar(30)" json:"user_agent"`
	Location  string     `gorm:"type:varchar(50)" json:"location"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

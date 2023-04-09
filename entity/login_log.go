package entity

import "time"

type LoginLog struct {
	UserID    int        `gorm:"index" json:"user_id"`
	Action    string     `gorm:"type:varchar(30)" json:"action"`
	UserAgent string     `gorm:"type:text" json:"user_agent"`
	Location  string     `gorm:"type:varchar(50)" json:"location"`
	OS        string     `gorm:"type:varchar(50)" json:"os"`
	Browser   string     `gorm:"type:varchar(50)" json:"browser"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

package entity

import "time"

type LoginLog struct {
	UserID    int        `gorm:"index" json:"user_id"`
	Action    string     `gorm:"type:varchar(30)" json:"action"`
	UserAgent string     `gorm:"type:varchar(30)" json:"user_agent"`
	Location  string     `gorm:"type:varchar(50)" json:"location"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

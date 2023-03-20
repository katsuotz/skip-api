package entity

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type Tabler interface {
	TableName() string
}

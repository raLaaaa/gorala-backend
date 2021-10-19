package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint `gorm:"primaryKey"`
	Description string
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

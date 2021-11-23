package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID            uint `gorm:"primaryKey"`
	Description   string
	ExecutionDate time.Time
	UserID        uint
	IsFinished    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

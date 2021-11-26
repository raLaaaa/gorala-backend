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
	IsFinished    bool `gorm:"default:false"`
	IsCarryOnTask bool `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

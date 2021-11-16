package models

import (
	"time"

	"gorm.io/gorm"
)

type ConfirmationToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex"`
	UserID    uint
	Activated bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

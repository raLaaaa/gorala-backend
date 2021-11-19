package models

import (
	"time"

	"gorm.io/gorm"
)

type SentMail struct {
	ID        uint `gorm:"primaryKey"`
	Recipient string
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

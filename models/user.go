package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                uint   `gorm:"primaryKey"`
	Email             string `gorm:"uniqueIndex"`
	Password          string
	Accepted          bool
	ConfirmationToken ConfirmationToken `gorm:"foreignKey:UserID"`
	AllTasks          []Task            `gorm:"foreignKey:UserID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

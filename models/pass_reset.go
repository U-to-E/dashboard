package models

import (
	"time"
)

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"not null"`
	Token     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

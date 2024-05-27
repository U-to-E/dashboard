package models

import "time"

type QuizState struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	QuizID    string
	Questions string `gorm:"type:text"` // JSON encoded questions
	StartTime time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

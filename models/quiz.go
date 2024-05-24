package models

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Name     string
	QuizID   string
	Duration string
	Path     string
}

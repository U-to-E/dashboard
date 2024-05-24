package models

import "gorm.io/gorm"

type Marks struct {
	gorm.Model
	QuizName         string
	CollageID        string
	MentorID         string
	QuizID           string
	StudentID        string
	TotalQuestions   string
	CorrectQuestions string
}

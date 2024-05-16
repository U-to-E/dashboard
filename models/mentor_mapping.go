package models

import "gorm.io/gorm"

type Mapping struct {
	gorm.Model
	MentorID   string  `gorm:"column:mentor_id;not null"`
	CollegeID  string  `gorm:"column:college_id;not null"`
	StudentsID []uint8 `gorm:"type:integer[]"`
}
